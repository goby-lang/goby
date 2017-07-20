package vm

import (
	"fmt"
	"plugin"
	"reflect"
)

// PluginObject is a special type that contains file pointer so we can keep track on target file.
type PluginObject struct {
	*baseObj
	fn     string
	plugin *plugin.Plugin
}

func (vm *VM) initPluginObject(fn string, p *plugin.Plugin) *PluginObject {
	return &PluginObject{fn: fn, plugin: p, baseObj: &baseObj{class: vm.topLevelClass(pluginClass)}}
}

func (vm *VM) initPluginClass() *RClass {
	pc := vm.initializeClass(pluginClass, false)
	pc.setBuiltInMethods(builtinPluginClassMethods(), true)
	pc.setBuiltInMethods(builtinPluginInstanceMethods(), false)
	vm.objectClass.setClassConstant(pc)
	return pc
}

// Only initialize file related methods after it's being required.
func builtinPluginClassMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{}
}

// Only initialize file related methods after it's being required.
func builtinPluginInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{
			Name: "send",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					s, ok := args[0].(*StringObject)

					if !ok {
						return t.vm.initErrorObject(TypeError, WrongArgumentTypeFormat, stringClass, args[0].Class().Name)
					}

					funcName := s.Value
					r := receiver.(*PluginObject)
					p := r.plugin
					f, err := p.Lookup(funcName)

					if err != nil {
						return t.vm.initErrorObject(InternalError, err.Error())
					}

					funcArgs := make([]reflect.Value, len(args)-1)

					for i, arg := range args[1:] {
						v, ok := arg.(builtInType)

						if ok {
							funcArgs[i] = reflect.ValueOf(v.value())
						} else {
							return t.vm.initErrorObject(InternalError, "Can't pass %s type object when calling go function", arg.Class().Name)
						}
					}

					fmt.Println(funcArgs)
					var ptr reflect.Value
					value := reflect.ValueOf(f)
					if value.Type().Kind() == reflect.Ptr {
						ptr = value
						value = ptr.Elem() // acquire value referenced by pointer
					} else {
						ptr = reflect.New(reflect.TypeOf(f)) // create new pointer
						temp := ptr.Elem()                   // create variable to value of pointer
						temp.Set(value)                      // set value of variable to our passed in value
					}

					result := reflect.ValueOf(reflect.ValueOf(f).Call(funcArgs)).Interface()

					switch result := result.(type) {
					case []reflect.Value:
						if len(result) == 1 {
							return t.vm.initStructObject(result[0])
						}

						structs := []Object{}
						for _, v := range result {
							structs = append(structs, t.vm.initStructObject(v))
						}

						return t.vm.initArrayObject(structs)
					default:
						return t.vm.initStructObject(result)
					}

				}
			},
		},
	}
}

// Polymorphic helper functions -----------------------------------------

// toString returns detailed infoof a array include elements it contains
func (p *PluginObject) toString() string {
	return "<Plugin: " + p.fn + ">"
}

// toJSON converts the receiver into JSON string.
func (f *PluginObject) toJSON() string {
	return f.toString()
}
