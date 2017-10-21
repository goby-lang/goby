package vm

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"unicode"

	"github.com/fatih/structs"
	"github.com/goby-lang/goby/vm/classes"
	"github.com/gorilla/mux"
)

type request struct {
	Method           string
	Body             string
	URL              string
	Path             string
	Host             string
	Protocol         string
	Headers          map[string][]string
	ContentLength    int64
	TransferEncoding []string
}

type response struct {
	status      int
	body        string
	contentType string
}

// Instance methods -----------------------------------------------------
func builtinSimpleServerInstanceMethods() []*BuiltinMethodObject {
	router := mux.NewRouter()

	return []*BuiltinMethodObject{
		{
			Name: "mount",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					path := args[0].(*StringObject).value
					method := args[1].(*StringObject).value

					router.HandleFunc(path, newHandler(t, blockFrame)).Methods(method)

					return receiver
				}
			},
		},
		{
			Name: "static",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					prefix := args[0].(*StringObject).value
					fileName := args[1].(*StringObject).value
					router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(http.Dir(fileName))))

					return receiver
				}
			},
		},
		{
			Name: "start",
			Fn: func(receiver Object, sourceLine int) builtinMethodBody {
				return func(t *thread, args []Object, blockFrame *normalCallFrame) Object {
					var port string
					var serveStatic bool
					server := receiver.(*RObject)

					portVar, ok := server.InstanceVariables.get("@port")

					if !ok {
						port = "8080"
					} else {
						port = portVar.(*StringObject).value
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

					router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						log.Printf("%s %s %s %d\n", r.Method, r.URL.Path, r.Proto, 404)
					})

					if serveStatic && fileRoot.Class() != t.vm.objectClass.getClassConstant(classes.NullClass) {
						fr := fileRoot.(*StringObject).value
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

// Internal functions ===================================================

// Functions for initialization -----------------------------------------

func initSimpleServerClass(vm *VM) {
	initHTTPClass(vm)
	net := vm.loadConstant("Net", true)
	simpleServer := vm.initializeClass("SimpleServer", false)
	simpleServer.setBuiltinMethods(builtinSimpleServerInstanceMethods(), false)
	net.setClassConstant(simpleServer)

	vm.execGobyLib("net/simple_server.gb")
}

// Other helper functions -----------------------------------------------

func newHandler(t *thread, blockFrame *normalCallFrame) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Go creates one goroutine per request, so we also need to create a new Goby thread for every request.
		thread := t.vm.newThread()
		res := httpResponseClass.initializeInstance()

		req := initRequest(t, w, r)
		result := thread.builtinMethodYield(blockFrame, req, res)

		if err, ok := result.Target.(*Error); ok {
			log.Printf("Error: %s", err.message)
			res.instanceVariableSet("@status", t.vm.initIntegerObject(500))
		}

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
	r.Protocol = req.Proto
	r.Headers = req.Header
	r.Body = string(body)
	r.ContentLength = req.ContentLength
	r.TransferEncoding = req.TransferEncoding
	r.Host = req.Host
	r.Path = req.URL.Path
	r.URL = req.RequestURI

	m := structs.Map(r)

	for k, v := range m {
		varName := "@" + toSnakeCase(k)
		reqObj.instanceVariableSet(varName, t.vm.initObjectFromGoType(v))
	}

	vars := map[string]Object{}

	for k, v := range mux.Vars(req) {
		vars[k] = t.vm.initStringObject(v)
	}

	reqObj.instanceVariableSet("@params", t.vm.initHashObject(vars))

	return reqObj
}

func setupResponse(w http.ResponseWriter, req *http.Request, res *RObject) {
	r := &response{}

	resStatus, ok := res.instanceVariableGet("@status")

	if ok {
		r.status = resStatus.(*IntegerObject).value
	} else {
		r.status = http.StatusOK
	}

	resBody, ok := res.instanceVariableGet("@body")

	if !ok {
		r.body = ""
	} else {
		r.body = resBody.(*StringObject).value
	}

	h, ok := res.instanceVariableGet("@headers")

	if headers, isHashObject := h.(*HashObject); ok && isHashObject {
		for k, v := range headers.Pairs {
			w.Header().Set(k, v.(*StringObject).value)
		}
	} else {
		r.contentType = "text/plain; charset=utf-8"
		w.Header().Set("Content-Type", r.contentType) // normal header
	}

	w.WriteHeader(r.status)

	io.WriteString(w, r.body)
	log.Printf("%s %s %s %d\n", req.Method, req.URL.Path, req.Proto, r.status)
}

func toSnakeCase(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
