package vm

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type response struct {
	status int
	body   string
}

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
				log.Fatal(http.ListenAndServe(":"+port, nil))
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
					req := initRequest(r)
					res := httpResponseClass.initializeInstance()

					v.builtInMethodYield(blockFrame, req, res)

					setupResponse(w, r, res)
				})

				return receiver
			}
		},
	},
}

func initRequest(r *http.Request) *RObject {
	req := httpRequestClass.initializeInstance()

	req.InstanceVariables.set("@method", initializeString(r.Method))
	req.InstanceVariables.set("@body", initializeString(""))
	req.InstanceVariables.set("@path", initializeString(r.URL.Path))
	req.InstanceVariables.set("@url", initializeString(r.URL.RequestURI()))

	return req
}

func setupResponse(w http.ResponseWriter, req *http.Request, res *RObject) {
	r := &response{}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8") // normal header

	resStatus, ok := res.InstanceVariables.get("@status")

	if ok {
		r.status = resStatus.(*IntegerObject).Value
	} else {
		r.status = http.StatusOK
	}

	resBody, ok := res.InstanceVariables.get("@body")

	if !ok {
		r.body = ""
	} else {
		r.body = resBody.(*StringObject).Value
	}

	io.WriteString(w, resBody.(*StringObject).Value)
	fmt.Printf("%s %s %s %d\n", req.Method, req.URL.Path, req.Proto, r.status)
}
