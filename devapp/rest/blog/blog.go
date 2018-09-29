package blog

import (
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/log"
)

type Blog struct {
	eel.RestResource
}

func (this *Blog) Resource() string {
	return "blog.blog"
}

func (this *Blog) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":    []string{"id:int"},
		"PUT":    []string{"type:int", "content", "info:json", "person:json-array", "?name"},
		"DELETE": []string{"id:int"},
	}
}

func (this *Blog) Get(ctx *eel.Context) {
	name := ctx.Request.GetString("name")
	ctx.Response.JSON(eel.Map{
		"name": name,
	})
}

func (this *Blog) Put(ctx *eel.Context) {
	content := ctx.Request.GetString("content")
	aType, _ := ctx.Request.GetInt("type")
	name := ctx.Request.GetString("name", "not_found")
	info := ctx.Request.GetJSON("info")
	persons := ctx.Request.GetJSONArray("person")
	log.Info("in blog.Put")
	ctx.Response.JSON(eel.Map{
		"type": aType,
		"content": content,
		"name": name,
		"info": info,
		"persons": persons,
	})
}

func (this *Blog) Post(ctx *eel.Context) {
	ctx.Response.JSON(eel.Map{
		"method": "post",
	})
}

func (this *Blog) Delete(ctx *eel.Context) {
	ctx.Response.JSON(eel.Map{
		"method": "delete",
	})
}