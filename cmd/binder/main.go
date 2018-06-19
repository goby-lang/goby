package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"

	. "github.com/dave/jennifer/jen"
)

var (
	in       = flag.String("in", ".", "folder to create bindings from")
	typeName = flag.String("type", "", "type to generate bindings for")
)

const (
	vmPkg = "github.com/goby-lang/goby/vm"
)

func typeFromExpr(e ast.Expr) string {
	var name string
	switch t := e.(type) {
	case *ast.Ident:
		name = t.Name

	case *ast.StarExpr:
		name = fmt.Sprintf("*%s", typeFromExpr(t.X))

	case *ast.SelectorExpr:
		name = fmt.Sprintf("%s.%s", typeFromExpr(t.X), t.Sel.Name)
	default:
		log.Printf("%T", e)

	}
	return name
}

func typeNameFromExpr(e ast.Expr) string {
	var name string
	switch t := e.(type) {
	case *ast.Ident:
		name = t.Name

	case *ast.StarExpr:
		name = typeFromExpr(t.X)

	case *ast.SelectorExpr:
		name = fmt.Sprintf("%s.%s", typeFromExpr(t.X), t.Sel.Name)
	default:
		log.Printf("%T", e)

	}
	return name
}

type Binding struct {
	ClassName       string
	ClassMethods    []*ast.FuncDecl
	InstanceMethods []*ast.FuncDecl
}

func (b *Binding) staticName() string {
	return fmt.Sprintf("_static%s", b.ClassName)
}

func (b *Binding) bindingName(f *ast.FuncDecl) string {
	return fmt.Sprintf("_binding_%s_%s", b.ClassName, f.Name.Name)
}

func (b *Binding) BindMethods(f *File) {
	f.Var().Id(b.staticName()).Op("=").New(Id(b.ClassName))
	for _, c := range b.ClassMethods {
		b.BindClassMethod(f, c)
		f.Line()
	}
	for _, c := range b.InstanceMethods {
		b.BindInstanceMethod(f, c)
		f.Line()
	}
}

func (b *Binding) BindClassMethod(f *File, d *ast.FuncDecl) {
	r := Id("r").Op(":=").Id(b.staticName()).Line()
	b.body(r, f, d)
}
func (b *Binding) BindInstanceMethod(f *File, d *ast.FuncDecl) {
	r := List(Id("r"), Id("ok")).Op(":=").Add(Id("receiver")).Assert(Op("*").Id(b.ClassName)).Line()
	r = r.If(Op("!").Id("ok")).Block(
		Panic(Lit("NOT OK")),
	).Line()
	b.body(r, f, d)
}

func (b *Binding) body(receiver *Statement, f *File, d *ast.FuncDecl) {
	s := f.Func().Id(b.bindingName(d))
	s = s.Params(Id("receiver").Qual(vmPkg, "Object"), Id("line").Id("int")).Qual(vmPkg, "Method")
	ff := Func().Params(Id("t").Op("*").Qual(vmPkg, "Thread"), Id("args").Index().Qual(vmPkg, "Object"))
	ff = ff.Qual(vmPkg, "Object")

	var args []*Statement
	for i, a := range d.Type.Params.List {
		if i == 0 {
			continue
		}
		i = i - 1
		c := List(Id(fmt.Sprintf("arg%d", i)), Id("ok")).Op(":=").Id("args").Index(Lit(i)).Assert(Id(typeFromExpr(a.Type)))
		c = c.Line()
		c = c.If(Op("!").Id("ok")).Block(
			Panic(Lit("NOT OK")),
		).Line()
		args = append(args, c)
	}

	inner := receiver.If(Len(Id("args")).Op("!=").Lit(d.Type.Params.NumFields() - 1)).Block(
		Panic(Lit("NOT OK")),
	).Line()

	argNames := []Code{
		Id("t"),
	}
	for i, a := range args {
		inner = inner.Add(a).Line()
		argNames = append(argNames, Id(fmt.Sprintf("arg%d", i)))
	}

	inner = inner.Return(Id("r").Dot(d.Name.Name).Call(argNames...))
	ff = ff.Block(inner)
	s.Block(Return(ff))

}

func main() {
	flag.Parse()

	fs := token.NewFileSet()
	buff, err := ioutil.ReadFile(*in)
	if err != nil {
		log.Fatal(err)
	}

	f, err := parser.ParseFile(fs, *in, string(buff), parser.AllErrors)
	if err != nil {
		log.Fatal(err)
	}

	bindings := make(map[string]*Binding)

	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Recv != nil {
				// class or instance?
				r := n.Recv.List[0]
				name := typeNameFromExpr(r.Type)

				b, ok := bindings[name]
				if !ok {
					b = new(Binding)
					b.ClassName = name
					bindings[name] = b
				}

				// class
				if r.Names == nil {
					b.ClassMethods = append(b.ClassMethods, n)
				} else {
					b.InstanceMethods = append(b.InstanceMethods, n)
				}
			}
		case *ast.TypeSpec:
			bindings[n.Name.Name] = &Binding{
				ClassName: n.Name.Name,
			}

		}

		return true
	})

	bnd, ok := bindings[*typeName]
	if !ok {
		log.Fatal("Uknown type", *typeName)
	}

	o := NewFile(f.Name.Name)
	bnd.BindMethods(o)

	err = o.Save("bindings.go")
	if err != nil {
		log.Fatal(err)
	}
}
