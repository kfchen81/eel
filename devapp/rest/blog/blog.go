package blog

import (
	"github.com/kfchen81/eel"
)

type Blog struct {
	eel.RestResource
}

func (this *Blog) Resource() string {
	return "blog.blog"
}

func (this *Blog) Get(ctx *eel.Context) {
	name := ctx.Request.GetString("name")
	ctx.Response.JSON(eel.Map{
		"name": name,
	})
}