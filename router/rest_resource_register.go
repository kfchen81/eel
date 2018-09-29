package router

import (
	"net/http"
	"github.com/kfchen81/eel/log"
	"sync"
	"github.com/kfchen81/eel/handler"
	"strings"
	"fmt"
	"time"
)

type RestResourceRegister struct {
	endpoint2resource map[string]RestResourceInterface
	pool sync.Pool
	sync.RWMutex
}

// Implement http.Handler interface.
func (this *RestResourceRegister) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	log.Infow("request ", "path", req.URL.Path, "method", req.Method)
	
	//get handler.Context from pool
	context := this.pool.Get().(*handler.Context)
	context.Reset(resp, req)
	defer this.pool.Put(context)
	
	endpoint := req.URL.Path
	if endpoint[len(endpoint)-1] != '/' {
		endpoint = endpoint + "/"
	}
	
	resource := this.endpoint2resource[endpoint]
	log.Debug(this.endpoint2resource)
	log.Debug(endpoint)
	if resource == nil {
		context.Response.ErrorWithCode(handler.HTTP_404, "resource:not_found", "无效的endpoint")
	} else {
		method := req.Method
		switch method {
		case "GET":
			resource.Get(context)
		}
	}
	
	timeDur := time.Since(startTime)
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
		gRegister.endpoint2resource = make(map[string]RestResourceInterface)
	}
	
	return gRegister
}

func RegisterResource(resource RestResourceInterface) {
	gRegister.Lock()
	defer gRegister.Unlock()
	
	endpoint := resource.Resource()
	pos := strings.LastIndex(endpoint, ".")
	log.Debug(pos)
	endpoint = fmt.Sprintf("/%s/%s/", endpoint[0:pos], endpoint[pos+1:])
	//apiEndpoint := fmt.Sprintf("/%s/%s/", endpoint[0:pos], endpoint[pos+1:])
	gRegister.endpoint2resource[endpoint] = resource
	log.Infow("register rest resource", "endpoint", endpoint)
	log.Debug(gRegister.endpoint2resource)
}

func init() {
	NewRestResourceRegister()
}