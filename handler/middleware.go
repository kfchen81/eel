package handler

type MiddlewareInterface interface {
	ProcessRequest(ctx *Context)
	ProcessResponse(ctx *Context)
}

type Middleware struct {

}

func (m *Middleware) ProcessRequest(ctx *Context) {

}

func (m *Middleware) ProcessResponse(ctx *Context) {

}


