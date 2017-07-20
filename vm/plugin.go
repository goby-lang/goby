package vm

import (
	"plugin"
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

					switch f := f.(type) {
					case func():
						f()
						return NULL
					case func(string) interface{}:
						arg := args[1].(*StringObject).Value
						return t.vm.initStructObject(f(arg))
					default:
						return FALSE
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
