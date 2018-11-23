package handler

import (
	"net/http"
	"context"
	"github.com/kfchen81/eel/log"
)

type RestResourceInterface interface {
	//Init(ct *handler.Context, controllerName, actionName string, app interface{})
	Prepare(ctx *Context)
	Get(ctx *Context)
	Post(ctx *Context)
	Delete(ctx *Context)
	Put(ctx *Context)
	//Finish(ctx *Context)
	Resource() string
	//EnableHTMLResource() bool
	GetParameters() map[string][]string
	//GetBusinessContext() context.Context
	//SeteelController(ctx *eel_context.Context, data map[interface{}]interface{})
}

/*RestResource 扩展eel.Controller, 作为rest中各个资源的基类
 */
type RestResource struct {
	Ctx  *context.Context
}

// Get adds a request function to handle GET request.
func (r *RestResource) Get(ctx *Context) {
	http.Error(ctx.Response.ResponseWriter, "Method Not Allowed", 405)
}

// Post adds a request function to handle POST request.
func (r *RestResource) Post(ctx *Context) {
	http.Error(ctx.Response.ResponseWriter, "Method Not Allowed", 405)
}

// Delete adds a request function to handle DELETE request.
func (r *RestResource) Delete(ctx *Context) {
	http.Error(ctx.Response.ResponseWriter, "Method Not Allowed", 405)
}

// Put adds a request function to handle PUT request.
func (r *RestResource) Put(ctx *Context) {
	http.Error(ctx.Response.ResponseWriter, "Method Not Allowed", 405)
}

// Head adds a request function to handle HEAD request.
func (r *RestResource) Head(ctx *Context) {
	http.Error(ctx.Response.ResponseWriter, "Method Not Allowed", 405)
}

// Patch adds a request function to handle PATCH request.
func (r *RestResource) Patch(ctx *Context) {
	http.Error(ctx.Response.ResponseWriter, "Method Not Allowed", 405)
}

// Options adds a request function to handle OPTIONS request.
func (r *RestResource) Options(ctx *Context) {
	http.Error(ctx.Response.ResponseWriter, "Method Not Allowed", 405)
}

func (r *RestResource) Resource() string {
	return ""
}

/*EnableHTMLResource 是否开启html资源
 */
func (r *RestResource) EnableHTMLResource() bool {
	return false
}

/*Parameters 获取需要检查的参数
 */
func (r *RestResource) GetParameters() map[string][]string {
	return nil
}

func (r *RestResource) GetBusinessContext() context.Context {
	return nil
	//data := r.Ctx.Input.GetData("bContext")
	//if data == nil {
	//	return nil
	//} else {
	//	bCtx := data.(context.Context)
	//	return bCtx
	//}
}

// prepare for request handling
func (r *RestResource) Prepare(ctx *Context) {
	log.Logger.Debug("in RestResource's Prepare...")
}

//func (r *RestResource) Finish() {
//	bCtx := r.GetBusinessContext()
//	if bCtx != nil {
//		o := bCtx.Value("orm")
//		if o != nil {
//			o.(orm.Ormer).Commit()
//			eel.Info("[ORM] commit transaction 2")
//		}
//
//		span := opentracing.SpanFromContext(bCtx)
//		if span != nil {
//			eel.Info("[Tracing] finish span in Controller.Finish")
//			span.(opentracing.Span).Finish()
//		}
//	}
//}
