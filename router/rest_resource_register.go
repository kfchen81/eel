package router

import (
	"net/http"
	"github.com/kfchen81/eel/log"
	"sync"
	"github.com/kfchen81/eel/handler"
	"strings"
	"fmt"
	"time"
	"reflect"
)

type RestResourceRegister struct {
	endpoint2resource map[string]handler.RestResourceInterface
	middlewares []handler.MiddlewareInterface
	pool sync.Pool
	sync.RWMutex
}

// Implement http.Handler interface.
func (this *RestResourceRegister) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	log.Infow("request ", "path", req.URL.Path, "method", req.Method)
	//
	//if AppPath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
	//	panic(err)
	//}
	//workPath, err := os.Getwd()
	//if err != nil {
	//	panic(err)
	//}
	
	//get handler.Context from pool
	context := this.pool.Get().(*handler.Context)
	context.Reset(resp, req)
	defer this.pool.Put(context)
	defer handler.RecoverPanic(context)
	
	//determine the resource will handle the request
	endpoint := req.URL.Path
	if endpoint[len(endpoint)-1] != '/' {
		endpoint = endpoint + "/"
	}
	resource := this.endpoint2resource[endpoint]
	if resource == nil {
		//resource is not exists
		context.Response.ErrorWithCode(http.StatusNotFound, "resource:not_found", "无效的endpoint", "")
	} else {
		//resource found, go through middlewares
		for _, middleware := range this.middlewares {
			middleware.ProcessRequest(context)
			
			if context.Response.Started {
				context.Response.Flush()
				goto FinishHandle
			}
		}
		
		//check resource params
		handler.CheckArgs(resource, context)
		if context.Response.Started {
			context.Response.Flush()
			goto FinishHandle
		}
		
		//pass param check, do prepare
		resource.Prepare(context)
		method := context.Request.Method()
		log.Infow("http request method", "method", method)
		log.Infow("http params", "params", context.Request.Input())
		switch method {
		case "GET":
			resource.Get(context)
		case "PUT":
			resource.Put(context)
		case "POST":
			resource.Post(context)
		case "DELETE":
			resource.Delete(context)
		default:
			http.Error(context.Response.ResponseWriter, "Method Not Allowed", 405)
		}
	}
	
	FinishHandle:
	timeDur := time.Since(startTime)
	context.Response.Body([]byte("robert"))
	context.Response.JSON(handler.Map{
		"name":"python",
	})
	statusCode := context.Response.Status
	contentLength := context.Response.ContentLength
	accessLog := fmt.Sprintf("%s - - [%s] \"%s %s %s %d %d\" %f", req.RemoteAddr, startTime.Format("02/Jan/2006 03:04:05"), req.Method, req.RequestURI, req.Proto, statusCode, contentLength, timeDur.Seconds())
	log.Infow(accessLog, "timeDur", timeDur.Seconds(), "status", statusCode)
}

// global resource register
var gRegister *RestResourceRegister = nil

// Create new RestResourceRegister as a http.Handler
func NewRestResourceRegister() *RestResourceRegister {
	if (gRegister == nil) {
		log.Debug("create global RestResourceRegister")
		gRegister = &RestResourceRegister{}
		gRegister.pool.New = func() interface{} {
			return handler.NewContext()
		}
		gRegister.endpoint2resource = make(map[string]handler.RestResourceInterface)
		gRegister.middlewares = make([]handler.MiddlewareInterface, 0)
	}
	
	return gRegister
}

func DoRegisterResource(resource handler.RestResourceInterface) {
	gRegister.Lock()
	defer gRegister.Unlock()
	
	endpoint := resource.Resource()
	pos := strings.LastIndex(endpoint, ".")
	endpoint = fmt.Sprintf("/%s/%s/", endpoint[0:pos], endpoint[pos+1:])
	//apiEndpoint := fmt.Sprintf("/%s/%s/", endpoint[0:pos], endpoint[pos+1:])
	gRegister.endpoint2resource[endpoint] = resource
	log.Infow("register rest resource", "endpoint", endpoint)
}

func DoRegisterMiddleware(middleware handler.MiddlewareInterface) {
	gRegister.Lock()
	defer gRegister.Unlock()
	
	gRegister.middlewares = append(gRegister.middlewares, middleware)
	log.Infow("register middleware", "name", reflect.TypeOf(middleware))
}

func init() {
	NewRestResourceRegister()
}