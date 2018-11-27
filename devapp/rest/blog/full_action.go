package blog

import (
	"github.com/kfchen81/eel"
)

type FullAction struct {
	eel.RestResource
}

func (this *FullAction) Resource() string {
	return "blog.full_action"
}

func (this *FullAction) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":    []string{},
	}
}

func (this *FullAction) Get(ctx *eel.Context) {
	//blogService := b_blog.NewBlogService(ctx.GetBusinessContext())
	//blogService.FullAction()
	ctx.Response.JSON(eel.Map{
	})
}
