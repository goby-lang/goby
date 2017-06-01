package vm

import (
	"fmt"
	"io"
	"net/http"
)

func initializeSimpleServerClass(vm *VM) {
	initializeHTTPClass(vm)
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
				serverClass := v.constants["Net"].returnClass().constants["SimpleServer"].returnClass()
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
				req := httpRequestClass.initializeInstance()
				res := httpResponseClass.initializeInstance()

				http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
					req.InstanceVariables.set("@method", initializeString(r.Method))
					req.InstanceVariables.set("@body", initializeString(""))
					req.InstanceVariables.set("@path", initializeString(r.URL.Path))
					req.InstanceVariables.set("@url", initializeString(r.URL.RequestURI()))

					builtInMethodYield(v, blockFrame, req, res)

					w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header

					resStatus, ok := res.InstanceVariables.get("@status")

					if ok {
						w.WriteHeader(resStatus.(*IntegerObject).Value)
					} else {
						w.WriteHeader(http.StatusOK)
					}

					resBody, ok := res.InstanceVariables.get("@body")

					if !ok {
						io.WriteString(w, "")
						return
					}

					io.WriteString(w, resBody.(*StringObject).Value)
				})

				return receiver
			}
		},
	},
}
