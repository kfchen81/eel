package handler

import (
	"net/http"
)

// NewContext return the Context with Input and Output
func NewContext() *Context {
	return &Context{
	}
}

type Context struct {
	Request  *Request
	Response *Response
	Data map[string]interface{}
}

// Reset init Context
func (ctx *Context) Reset(rw http.ResponseWriter, r *http.Request) {
	if ctx.Request == nil {
		ctx.Request = &Request{}
	}
	ctx.Request.Reset(r)
	
	if ctx.Response == nil {
		ctx.Response = &Response{}
	}
	ctx.Response.Reset(rw)
}

func (ctx *Context) SetJSON(key string, value map[string]interface{}) {
	ctx.Request.SetJSON(key, value)
}

func (ctx *Context) SetJSONArray(key string, value []interface{}) {
	ctx.Request.SetJSONArray(key, value)
}