package vm

import (
	"fmt"
	"github.com/st0012/metago"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
)

func (vm *VM) initPluginObject(fn string, p *plugin.Plugin) *PluginObject {
	return &PluginObject{fn: fn, plugin: p, baseObj: &baseObj{class: vm.topLevelClass(pluginClass)}}
}

func initPluginClass(vm *VM) {
	pc := vm.initializeClass(pluginClass, false)
	pc.setBuiltInMethods(builtinPluginClassMethods(), true)
	pc.setBuiltInMethods(builtinPluginInstanceMethods(), false)
	vm.objectClass.setClassConstant(pc)

	vm.execGobyLib("plugin.gb")
}

type pluginContext struct {
	pkgs  []*pkg
	funcs []*function
}

func (c *pluginContext) importPkg(prefix, name string) {
	c.pkgs = append(c.pkgs, &pkg{Prefix: prefix, Name: name})
}

func (c *pluginContext) addFunc(prefix, name string) {
	c.funcs = append(c.funcs, &function{Prefix: prefix, Name: name})
}

// PluginObject is a special type that contains a Go's plugin
type PluginObject struct {
	*baseObj
	fn     string
	plugin *plugin.Plugin
}

// Polymorphic helper functions -----------------------------------------
func (p *PluginObject) toString() string {
	return "<Plugin: " + p.fn + ">"
}

func (p *PluginObject) toJSON() string {
	return p.toString()
}

func setPluginContext(context Object) *pluginContext {
	pc := &pluginContext{pkgs: []*pkg{}, funcs: []*function{}}

	funcs, _ := context.instanceVariableGet("@funcs")
	pkgs, _ := context.instanceVariableGet("@pkgs")

	fs := funcs.(*ArrayObject)
	ps := pkgs.(*ArrayObject)

	for _, f := range fs.Elements {
		fInfos := f.(*ArrayObject)
		prefix := fInfos.Elements[0].(*StringObject).value
		name := fInfos.Elements[1].(*StringObject).value

		pc.addFunc(prefix, name)
	}

	for _, p := range ps.Elements {
		pInfos := p.(*ArrayObject)
		prefix := pInfos.Elements[0].(*StringObject).value
		name := pInfos.Elements[1].(*StringObject).value

		pc.importPkg(prefix, name)
	}

	return pc
}

func builtinPluginClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}

func builtinPluginInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "compile",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					context, ok := receiver.instanceVariableGet("@context")

					if !ok {
						return NULL
					}

					pc := setPluginContext(context)

					fmt.Println(pc)
					pkgPath := args[0].(*StringObject).value
					goPath := os.Getenv("GOPATH")
					// This is to prevent some path like GODEP_PATH:GOPATH
					// which can happen on Travis CI
					ps := strings.Split(goPath, ":")
					goPath = ps[len(ps)-1]

					fullPath := filepath.Join(goPath, "src", pkgPath)
					_, pkgName := filepath.Split(fullPath)
					pkgName = strings.Split(pkgName, ".")[0]
					soName := filepath.Join("./", pkgName+".so")

					// Open plugin first
					p, err := plugin.Open(soName)

					// If there's any issue open a plugin, assume it's not well compiled
					if err != nil {
						cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", fmt.Sprintf("./%s.so", pkgName), fullPath)
						out, err := cmd.CombinedOutput()

						if err != nil {
							return t.vm.initErrorObject(InternalError, "Error: %s from %s", string(out), strings.Join(cmd.Args, " "))
						}

						p, err = plugin.Open(soName)

						if err != nil {
							return t.vm.initErrorObject(InternalError, "Error occurs when open %s package: %s", soName, err.Error())
						}
					}

					return t.vm.initPluginObject(fullPath, p)
				}
			},
		},
		{
			Name: "send",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, stringClass, args[0].Class().Name)
					}

					funcName := s.value
					r := receiver.(*PluginObject)
					p := r.plugin
					f, err := p.Lookup(funcName)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					funcArgs, err := convertToGoFuncArgs(args[1:])

					if err != nil {
						t.vm.initErrorObject(TypeError, err.Error())
					}

					funcValue := reflect.ValueOf(f)

					// Check if f is a pointer to function instead of function object
					if funcValue.Type().Kind() == reflect.Ptr {
						ptr := funcValue
						funcValue = ptr.Elem()
					}

					result := reflect.ValueOf(funcValue.Call(metago.WrapArguments(funcArgs...))).Interface()

					return t.vm.initObjectFromGoType(metago.UnwrapReflectValues(result))
				}
			},
		},
	}
}
