package vm

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"

	"github.com/goby-lang/goby/vm/classes"
	"github.com/goby-lang/goby/vm/errors"
	"github.com/st0012/metago"
)

// PluginObject is a special type that contains a Go's plugin
type PluginObject struct {
	*baseObj
	fn     string
	plugin *plugin.Plugin
}

// Class methods --------------------------------------------------------
func builtinPluginClassMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "new",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					if len(args) != 1 {
						return t.vm.initErrorObject(errors.ArgumentError, errors.WrongNumberOfArgumentFormat, 1, len(args))
					}

					name, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, "String", args[0].Class().Name)
					}

					return &PluginObject{fn: name.value, baseObj: &baseObj{class: t.vm.topLevelClass(classes.PluginClass), InstanceVariables: newEnvironment()}}
				}
			},
		},
		{
			Name: "use",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					pkgPath := args[0].(*StringObject).value
					_, pkgName := filepath.Split(pkgPath)
					pkgName = strings.Split(pkgName, ".")[0]
					soName := filepath.Join("./", pkgName+".so")

					p, err := compileAndOpenPlugin(soName, pkgPath)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					return t.vm.initPluginObject(pkgPath, p)
				}
			},
		},
	}
}

// Instance methods -----------------------------------------------------
func builtinPluginInstanceMethods() []*BuiltinMethodObject {
	return []*BuiltinMethodObject{
		{
			Name: "compile",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					r := receiver.(*PluginObject)
					context, ok := receiver.instanceVariableGet("@context")

					if !ok {
						return NULL
					}

					// Create plugins directory
					pluginDir := "./plugins"

					ok, err := fileExists(pluginDir)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					if !ok {
						os.Mkdir(pluginDir, 0777)
					}

					// generate plugin content from context
					pc := setPluginContext(context)
					pluginContent := compilePluginTemplate(pc.pkgs, pc.funcs)

					// create plugin file
					fn := fmt.Sprintf("%s/%s", pluginDir, r.fn)

					file, err := os.OpenFile(fn+".go", os.O_RDWR|os.O_CREATE, 0755)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, "Error when creating plugin: %s", err.Error())
					}

					file.WriteString(pluginContent)

					soName := fn + ".so"

					p, err := compileAndOpenPlugin(soName, file.Name())

					if err != nil {
						t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					r.plugin = p

					return r
				}
			},
		},
		{
			Name: "go_func",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(errors.TypeError, errors.WrongArgumentTypeFormat, classes.StringClass, args[0].Class().Name)
					}

					funcName := s.value
					r := receiver.(*PluginObject)
					p := r.plugin
					f, err := p.Lookup(funcName)

					if err != nil {
						return t.vm.initErrorObject(errors.InternalError, err.Error())
					}

					funcArgs, err := convertToGoFuncArgs(args[1:])

					if err != nil {
						t.vm.initErrorObject(errors.TypeError, err.Error())
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

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func (vm *VM) initPluginObject(fn string, p *plugin.Plugin) *PluginObject {
	return &PluginObject{fn: fn, plugin: p, baseObj: &baseObj{class: vm.topLevelClass(classes.PluginClass)}}
}

func initPluginClass(vm *VM) {
	pc := vm.initializeClass(classes.PluginClass, false)
	pc.setBuiltinMethods(builtinPluginClassMethods(), true)
	pc.setBuiltinMethods(builtinPluginInstanceMethods(), false)
	vm.objectClass.setClassConstant(pc)

	vm.execGobyLib("plugin.gb")
}

// Polymorphic helper functions -----------------------------------------

// `toString` returns the object's name as the string format
func (p *PluginObject) toString() string {
	return "<Plugin: " + p.fn + ">"
}

// `toJSON` just delegates to `toString`
func (p *PluginObject) toJSON() string {
	return p.toString()
}

// Other helper functions -----------------------------------------------

func setPluginContext(context Object) *pluginContext {
	pc := &pluginContext{pkgs: []*pkg{}, funcs: []*function{}}

	funcs, _ := context.instanceVariableGet("@functions")
	pkgs, _ := context.instanceVariableGet("@packages")

	fs := funcs.(*ArrayObject)
	ps := pkgs.(*ArrayObject)

	for _, f := range fs.Elements {
		fInfos := f.(*HashObject)
		prefix := fInfos.Pairs["prefix"].(*StringObject).value
		name := fInfos.Pairs["name"].(*StringObject).value

		pc.addFunc(prefix, name)
	}

	for _, p := range ps.Elements {
		pInfos := p.(*HashObject)
		prefix := pInfos.Pairs["prefix"].(*StringObject).value
		name := pInfos.Pairs["name"].(*StringObject).value

		pc.importPkg(prefix, name)
	}

	return pc
}

// `fileExists` returns whether the given file or directory exists or not
func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func compileAndOpenPlugin(soName, fileName string) (*plugin.Plugin, error) {
	// Open plugin first
	p, err := plugin.Open(soName)

	// If there's any issue open a plugin, assume it's not well compiled
	if err != nil {
		cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", soName, fileName)
		out, err := cmd.CombinedOutput()

		if err != nil {
			return nil, fmt.Errorf("Error: %s from %s", string(out), strings.Join(cmd.Args, " "))
		}

		p, err = plugin.Open(soName)

		if err != nil {
			return nil, fmt.Errorf("Error occurs when open %s package: %s", soName, err.Error())
		}
	}

	return p, nil
}

// Plugin context =======================================================

type pluginContext struct {
	pkgs  []*pkg
	funcs []*function
}

// Polymorphic helper functions -----------------------------------------

func (c *pluginContext) importPkg(prefix, name string) {
	c.pkgs = append(c.pkgs, &pkg{Prefix: prefix, Name: name})
}

func (c *pluginContext) addFunc(prefix, name string) {
	c.funcs = append(c.funcs, &function{Prefix: prefix, Name: name})
}
