package vm

import (
	"github.com/fatih/structs"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
)

type response struct {
	status      int
	body        string
	contentType string
}

type request struct {
	Method string
	Body   string
	URL    string
	Path   string
	Host   string
}

func initSimpleServerClass(vm *VM) {
	initHTTPClass(vm)
	net := vm.loadConstant("Net", true)
	simpleServer := vm.initializeClass("SimpleServer", false)
	simpleServer.setBuiltInMethods(builtinSimpleServerInstanceMethods(), false)
	net.setClassConstant(simpleServer)

	vm.execGobyLib("net/simple_server.gb")
}

func builtinSimpleServerInstanceMethods() []*BuiltInMethodObject {
	router := mux.NewRouter()

	return []*BuiltInMethodObject{
		{
			Name: "mount",
			Fn: func(receiver Object) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *callFrame) Object {
					path := args[0].(*StringObject).Value
					method := args[1].(*StringObject).Value
					router.HandleFunc(path, newHandler(t, blockFrame)).Methods(method)

					return receiver
				}
			},
		},
		{
			Name: "start",
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

					log.Println("SimpleServer start listening on port: " + port)

					c := make(chan os.Signal, 1)
					signal.Notify(c, os.Interrupt)

					go func() {
						for range c {
							log.Println("SimpleServer gracefully stopped")
							os.Exit(0)
						}
					}()

					fileRoot, serveStatic := server.InstanceVariables.get("@file_root")

					if serveStatic {
						fr := fileRoot.(*StringObject).Value
						currentDir, _ := os.Getwd()
						fp := filepath.Join(currentDir, fr)
						fs := http.FileServer(http.Dir(fp))
						http.Handle("/", fs)
					} else {
						http.Handle("/", router)
					}

					err := http.ListenAndServe(":"+port, nil)

					if err != http.ErrServerClosed { // HL
						log.Fatalf("listen: %s\n", err)
					}

					return receiver
				}
			},
		},
	}
}

func newHandler(t *thread, blockFrame *callFrame) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Go creates one goroutine per request, so we also need to create a new Goby thread for every request.
		thread := t.vm.newThread()
		res := httpResponseClass.initializeInstance()
		req := initRequest(t, w, r)
		thread.builtInMethodYield(blockFrame, req, res)
		thread = nil
		setupResponse(w, r, res)
	}
}

func initRequest(t *thread, w http.ResponseWriter, req *http.Request) *RObject {
	r := request{}
	reqObj := httpRequestClass.initializeInstance()

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		http.Error(w, "can't read body", http.StatusBadRequest)
		return reqObj
	}

	r.Method = req.Method
	r.Body = string(body)
	r.Path = req.URL.Path
	r.URL = req.RequestURI
	r.Host = req.Host

	m := structs.Map(r)

	for k, v := range m {
		varName := "@" + strings.ToLower(k)
		reqObj.InstanceVariables.set(varName, t.vm.initObjectFromGoType(v))
	}

	return reqObj
}

func setupResponse(w http.ResponseWriter, req *http.Request, res *RObject) {
	r := &response{}

	resStatus, ok := res.instanceVariableGet("@status")

	if ok {
		r.status = resStatus.(*IntegerObject).Value
	} else {
		r.status = http.StatusOK
	}

	resBody, ok := res.instanceVariableGet("@body")

	if !ok {
		r.body = ""
	} else {
		r.body = resBody.(*StringObject).Value
	}

	contentType, ok := res.instanceVariableGet("@content_type")

	if !ok {
		r.contentType = "text/plain; charset=utf-8"
	} else {
		r.contentType = contentType.toString()
	}

	w.WriteHeader(r.status)
	w.Header().Set("Content-Type", r.contentType) // normal header

	io.WriteString(w, r.body)
	log.Printf("%s %s %s %d\n", req.Method, req.URL.Path, req.Proto, r.status)
}
