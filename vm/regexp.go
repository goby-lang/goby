package vm

import (
	"fmt"

	"github.com/dlclark/regexp2"
)

var (
	regexpClass *RClass
)

// RegexpObject wraps github.com/dlclark/regexp2 library.
type RegexpObject struct {
	Class   *RClass
	Regexp2 regexp2.Regexp
}

func (regex *RegexpObject) toString() string {
	return fmt.Sprintf("/%s/", regex.Regexp2.String())
}

func (regex *RegexpObject) toJSON() string {
	return regex.toString()
}

func (regex *RegexpObject) returnClass() Class {
	return regex.Class
}

func initRegexpObject(value string) *RegexpObject {
	return &RegexpObject{
		Class:   regexpClass,
		Regexp2: regexp2.Regexp{},
	}
}

func initRegexpClass() {
	bc := &BaseClass{Name: "Regexp", ClassMethods: newEnvironment(), Methods: newEnvironment(), Class: classClass, pseudoSuperClass: objectClass, superClass: objectClass}
	regc := &RClass{BaseClass: bc}
	regc.setBuiltInMethods(builtInRegexpInstanceMethods, false)
	regexpClass = regc
}

var builtInRegexpInstanceMethods = []*BuiltInMethodObject{
	{
		Name: "compile",
		Fn: func(receiver Object) builtinMethodBody {
			return func(t *thread, args []Object, blockFrame *callFrame) Object {
				return initRegexpObject(args[0].toString())
			}
		},
	},
}
