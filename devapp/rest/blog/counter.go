package blog

import (
	b_blog "github.com/kfchen81/eel/devapp/business/blog"

	"github.com/kfchen81/eel"
)

type Counter struct {
	eel.RestResource
}

func (this *Counter) Resource() string {
	return "blog.counter"
}

func (this *Counter) GetParameters() map[string][]string {
	return map[string][]string{
		"POST":    []string{"id:int", "delta:int"},
	}
}

func (this *Counter) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	delta, _ := req.GetInt("delta")

	bCtx := ctx.GetBusinessContext()
	blogRepository := b_blog.NewBlogRepository(bCtx)
	blog := blogRepository.GetBlog(id)
	blog.UpdateCounter(delta)
	
	ctx.Response.JSON(eel.Map{})
}
