package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"

	"github.com/dave/jennifer/jen"
)

var (
	in  = flag.String("in", ".", "folder to create bindings from")
	out = flag.String("out", "bindings", "output file name")
)

const (
	vmPkg = "github.com/goby-lang/goby/vm"
)

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

func (b *Binding) BindMethods(f *jen.File) {
	for _, c := range b.ClassMethods {
		b.BindClassMethod(f, c)
		f.Line()
	}
	for _, c := range b.InstanceMethods {
		b.BindInstanceMethod(f, c)
		f.Line()
	}
}

func (b *Binding) BindClassMethod(f *jen.File, d *ast.FuncDecl) {
	f.Func().Id(b.bindingName(d)).Call().Block()
}

func (b *Binding) BindInstanceMethod(f *jen.File, d *ast.FuncDecl) {
	s := f.Func().Id(b.bindingName(d))
	s = s.Params(jen.Id("r").Qual(vmPkg, "Object"), jen.Id("line").Id("int")).Qual(vmPkg, "Method")
	ff := jen.Func().Params(jen.Id("t").Qual(vmPkg, "Thread"), jen.Id("args").Index().Qual(vmPkg, "Object"))
	ff = ff.Qual(vmPkg, "Object")
	ff = ff.Block(
		jen.Return(nil),
	)
	s.Block(jen.Return(ff))

	// func closeDB(receiver vm.Object, sourceLine int) vm.Method {
	// 	return func(t *vm.Thread, args []vm.Object) vm.Object {
	// 	}
	// 	return nil
	// }
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
				var name string
				switch t := r.Type.(type) {
				case *ast.Ident:
					name = t.Name

				case *ast.StarExpr:
					name = t.X.(*ast.Ident).Name
				}

				b, ok := bindings[name]
				if !ok {
					b := new(Binding)
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
	o := jen.NewFile(*out)

	for _, x := range bindings {
		x.BindMethods(o)
	}

	var b bytes.Buffer
	err = o.Render(&b)
	log.Println(err, b.String())
}
