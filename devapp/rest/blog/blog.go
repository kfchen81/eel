package blog

import (
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/devapp/business/blog"
	"github.com/kfchen81/eel/devapp/business/account"
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
		"PUT":    []string{"title", "content"},
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
	title := ctx.Request.GetString("title")
	content := ctx.Request.GetString("content")
	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	
	newBlog := blog.NewBlog(bCtx, user, title, content)

	eel.Logger.Info("in blog.Put")
	ctx.Response.JSON(eel.Map{
		"id": newBlog.Id,
		"title": newBlog.Title,
		"content": newBlog.Content,
	})
}

func (this *Blog) Post(ctx *eel.Context) {
	account := ctx.Get("account").(string)
	ctx.Response.JSON(eel.Map{
		"method": "post",
		"account": account,
	})
}

func (this *Blog) Delete(ctx *eel.Context) {
	ctx.Response.JSON(eel.Map{
		"method": "delete",
	})
}