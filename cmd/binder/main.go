package main

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
)

var (
	in  = flag.String("in", ".", "folder to create bindings from")
	out = flag.String("out", "bindings", "output file name")
)

type Binding struct {
	ClassName       string
	ClassMethods    []*ast.FuncDecl
	InstanceMethods []*ast.FuncDecl
}

func (Binding) New() {

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

	var binding Binding

	ast.Inspect(f, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.FuncDecl:
			if n.Recv != nil {
				// class or instance?
				r := n.Recv.List[0]
				r.Type.

				// class
				if r.Names == nil {

				} else {
					// instance
				}
				if n.Recv.List[0].Type {
					binding.ClassMethods = append(binding.ClassMethods, n)
				}
			}
		case *ast.TypeSpec:
			binding.ClassName = n.Name.Name

		}

		return true
	})
	log.Printf("%+v", binding)
	// o := jen.NewFile(*out)
	//
	// o.Func().Id("main").Params().Block(jen.Qual("log", "Println").Call(jen.Lit("Hello world")))
	// b := new(bytes.Buffer)
	// err = o.Render(b)
	// log.Println(err, b.String())
}
