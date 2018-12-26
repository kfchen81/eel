package router

import (
	"github.com/kfchen81/eel/handler"
	"sync"
)

import (
	go_context "context"
	"net/http"
	"github.com/kfchen81/eel/log"
	"strings"
	"fmt"
	"time"
	"reflect"
	"github.com/kfchen81/eel/config/env"
	"path/filepath"
	"github.com/kfchen81/eel/utils"
	"os"
	"io"
	"bytes"
	"sort"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/kfchen81/eel/tracing"
	"github.com/kfchen81/eel/config"
	"github.com/astaxie/beego"
	"net/http/pprof"
)

type RestResourceRegister struct {
	endpoint2resource map[string]handler.RestResourceInterface
	middlewares []handler.MiddlewareInterface
	pool sync.Pool
	sync.RWMutex
}

func ServeStaticFile(path string, response *handler.Response) bool {
	gopaths := env.Get("GOPATH", "")
	if gopaths == "" {
		return false
	}
	
	items := strings.Split(gopaths, ":")
	for _, item := range items {
		absPath := filepath.Join(item, "src/github.com/kfchen81/eel", path)
		isExists := utils.FileExists(absPath)
		log.Logger.Infow("check static file", "path", absPath, "exists", isExists)
		if isExists {
			file, err := os.Open(absPath)
			if err != nil {
				return false
			}
			defer file.Close()
			
			var bufferWriter bytes.Buffer
			io.Copy(&bufferWriter, file)
			
			contentType := "text/html; charset=utf-8"
			if strings.HasSuffix(path, ".css") {
				contentType = "text/css"
			} else if strings.HasSuffix(path, ".js") {
				contentType = "text/javascript"
			}
			response.Header("Content-Type", contentType)
			response.Body(bufferWriter.Bytes())
			
			return true
		} else {
		
		}
	}
	
	//AppPath, _ := filepath.Abs(filepath.Dir(os.Args[0]));
	//log.Infow("request ", "AppPath", AppPath)
	
	return false
}


// Implement http.Handler interface.
func (this *RestResourceRegister) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	log.Logger.Infow("request ", "path", req.URL.Path, "method", req.Method)
	
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
	if strings.HasPrefix(endpoint, "/debug/pprof/") {
		if endpoint == "/debug/pprof/profile" {
			pprof.Profile(resp, req)
		} else if endpoint == "/debug/pprof/cmdline" {
			pprof.Cmdline(resp, req)
		} else if endpoint == "/debug/pprof/symbol" {
			pprof.Symbol(resp, req)
		} else {
			pprof.Index(resp, req)
		}
		return
	}
	
	beego.Warn(endpoint)
	if endpoint[len(endpoint)-1] != '/' {
		endpoint = endpoint + "/"
	}
	resource := this.endpoint2resource[endpoint]
	if resource == nil {
		//resource is not exists
		ServeStaticFile(req.URL.Path, context.Response)
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
		
		//add tracing span
		spanCtx, _ := tracing.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
		operationName := fmt.Sprintf("%s %s", req.Method, endpoint)
		span := tracing.Tracer.StartSpan(operationName, ext.RPCServerOption(spanCtx))
		bCtx := context.GetBusinessContext()
		bCtx = opentracing.ContextWithSpan(bCtx, span)
		context.Set("rootSpan", span)
		//在结束时，report span
		defer func() {
			span := opentracing.SpanFromContext(bCtx)
			if span != nil {
				log.Logger.Debug("[Tracing] finish span in ServeHTTP")
				span.(opentracing.Span).Finish()
			}
		}()
		
		//add gorm's Transaction
		if config.Runtime.DB != nil {
			subSpan := tracing.CreateSubSpan(span, "db-begin")
			tx := config.Runtime.DB.Begin()
			subSpan.Finish()
			tx.InstantSet("rootSpan", span)
			bCtx = go_context.WithValue(bCtx, "orm", tx)
			context.Set("orm", tx)
		}
		
		context.SetBusinessContext(bCtx)
		
		//check resource params
		handler.CheckArgs(resource, context)
		if context.Response.Started {
			context.Response.Flush()
			goto FinishHandle
		}
		
		//pass param check, do prepare
		resource.Prepare(context)
		method := context.Request.Method()
		log.Logger.Infow("http request method", "method", method)
		log.Logger.Infow("http params", "params", context.Request.Input())
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
	log.Logger.Infow(accessLog, "timeDur", timeDur.Seconds(), "status", statusCode)
}

// global resource register
var gRegister *RestResourceRegister = nil

// Create new RestResourceRegister as a http.Handler
func NewRestResourceRegister() *RestResourceRegister {
	if (gRegister == nil) {
		log.Logger.Debug("create global RestResourceRegister")
		gRegister = &RestResourceRegister{}
		gRegister.pool.New = func() interface{} {
			return handler.NewContext()
		}
		gRegister.endpoint2resource = make(map[string]handler.RestResourceInterface)
		gRegister.middlewares = make([]handler.MiddlewareInterface, 0)
	}
	
	return gRegister
}


// Resources: get all registered resources
func Resources() []string {
	resources := make([]string, 0)
	if gRegister != nil {
		gRegister.Lock()
		defer gRegister.Unlock()
		
		for _, v := range gRegister.endpoint2resource {
			resource := v.(handler.RestResourceInterface).Resource()
			if resource != "console.console" {
				resources = append(resources, resource)
			}
		}
		
		sort.Strings(resources)
	}
	
	return resources
}

func DoRegisterResource(resource handler.RestResourceInterface) {
	gRegister.Lock()
	defer gRegister.Unlock()
	
	endpoint := resource.Resource()
	pos := strings.LastIndex(endpoint, ".")
	endpoint = fmt.Sprintf("/%s/%s/", endpoint[0:pos], endpoint[pos+1:])
	//apiEndpoint := fmt.Sprintf("/%s/%s/", endpoint[0:pos], endpoint[pos+1:])
	gRegister.endpoint2resource[endpoint] = resource
	log.Logger.Infow("[rest_resource] register rest resource", "endpoint", endpoint)
}

func DoRegisterMiddleware(middleware handler.MiddlewareInterface) {
	gRegister.Lock()
	defer gRegister.Unlock()
	
	gRegister.middlewares = append(gRegister.middlewares, middleware)
	log.Logger.Infow("[middleware] register middleware", "name", reflect.TypeOf(middleware))
}

func init() {
	NewRestResourceRegister()
}