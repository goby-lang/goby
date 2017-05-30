package vm

import (
	"fmt"
	"net/http"
)

func initializeSimpleServerClass(vm *VM) {
	net := vm.loadConstant("Net", true)
	simpleServer := initializeClass("SimpleServer", false)
	simpleServer.setBuiltInMethods(builtinSimpleServerClassMethods, true)
	simpleServer.setBuiltInMethods(builtinSimpleServerInstanceMethods, false)
	net.constants[simpleServer.Name] = &Pointer{simpleServer}
}

var builtinSimpleServerClassMethods = []*BuiltInMethodObject{
	{
		Name: "new",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				serverClass := v.constants["Net"].Target.(*RClass).constants["SimpleServer"].Target.(*RClass)
				server := serverClass.initializeInstance()
				server.InstanceVariables.set("@port", args[0])
				return server
			}
		},
	},
}

var builtinSimpleServerInstanceMethods = []*BuiltInMethodObject{
	{
		Name: "start",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				var port string

				portVar, ok := receiver.(*RObject).InstanceVariables.get("@port")

				if !ok {
					port = "8080"
				} else {
					port = portVar.(*StringObject).Value
				}

				fmt.Println("Start listening on port: " + port)
				http.ListenAndServe(":"+port, nil)
				return receiver
			}
		},
	},
	{
		Name: "mount",
		Fn: func(receiver Object) builtinMethodBody {
			return func(v *VM, args []Object, blockFrame *callFrame) Object {
				path := args[0].(*StringObject).Value

				http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					// args here should be request and response, which haven't been implemented yet.
					string := builtInMethodYield(v, blockFrame, nil).Target.(*StringObject)
					fmt.Fprint(w, string.Value)
				})

				return receiver
			}
		},
	},
}
