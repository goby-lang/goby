package vm

import (
	"github.com/fatih/structs"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"path/filepath"
)

type response struct {
	status int
	body   string
}

type request struct {
	Method string
	Body   string
	URL    string
	Path   string
}

func initializeSimpleServerClass(vm *VM) {
	initializeHTTPClass(vm)
	net := vm.loadConstant("Net", true)
	simpleServer := initializeClass("SimpleServer", false)
	simpleServer.setBuiltInMethods(builtinSimpleServerInstanceMethods(), false)
	net.constants[simpleServer.Name] = &Pointer{simpleServer}

	vm.execGobyLib("net/simple_server.gb")
}

func builtinSimpleServerInstanceMethods() []*BuiltInMethodObject {
	return []*BuiltInMethodObject{
		{Name: "start",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					var port string
					var serveStatic bool
					server := receiver.(*RObject)

					portVar, ok := server.InstanceVariables.get("@port")

					if !ok {
						port = "8080"
					} else {
						port = portVar.(*StringObject).Value
					}

					fileRoot, serveStatic := server.InstanceVariables.get("@file_root")

					if serveStatic {
						fr := fileRoot.(*StringObject).Value
						currentDir, _ := os.Getwd()
						fp := filepath.Join(currentDir, fr)
						fs := http.FileServer(http.Dir(fp))
						http.Handle("/", fs)
					}

					log.Println("SimpleServer start listening on port: " + port)

					c := make(chan os.Signal, 1)
					signal.Notify(c, os.Interrupt)

					go func() {
						for range c {
							log.Println("SimpleServer gracefully stopped")
							os.Exit(0)
						}
					}()

					err := http.ListenAndServe(":"+port, nil)

					if err != http.ErrServerClosed { // HL
						log.Fatalf("listen: %s\n", err)
					}

					return receiver
				}
			},
		},
		{
			Name: "mount",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					path := args[0].(*StringObject).Value

					http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
						// Go creates one goroutine per request, so we also need to create a new Goby thread for every request.
						thread := t.vm.newThread()
						res := httpResponseClass.initializeInstance()
						req := initRequest(r)
						thread.builtInMethodYield(blockFrame, req, res)
						thread = nil
						setupResponse(w, r, res)
					})

					return receiver
				}
			},
		},
	}
}

func initRequest(req *http.Request) *RObject {
	r := request{}
	reqObj := httpRequestClass.initializeInstance()

	r.Method = req.Method
	r.Body = ""
	r.Path = req.URL.Path
	r.URL = req.Host + req.RequestURI

	m := structs.Map(r)

	for k, v := range m {
		varName := "@" + strings.ToLower(k)
		reqObj.InstanceVariables.set(varName, initObject(v))
	}

	return reqObj
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

	io.WriteString(w, r.body)
	log.Printf("%s %s %s %d\n", req.Method, req.URL.Path, req.Proto, r.status)
}

func initObject(v interface{}) Object {
	switch v := v.(type) {
	case string:
		return initializeString(v)
	case int:
		return initilaizeInteger(v)
	case bool:
		if v {
			return TRUE
		}

		return FALSE
	default:
		panic("Can't init object")
	}
}
