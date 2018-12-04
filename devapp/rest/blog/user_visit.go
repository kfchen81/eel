package blog

import (
	"github.com/kfchen81/eel/devapp/business/blog"
	"github.com/kfchen81/eel/devapp/business/account"
	
	"github.com/kfchen81/eel"
)

type UserVisit struct {
	eel.RestResource
}

func (this *UserVisit) Resource() string {
	return "blog.user_visit"
}

func (this *UserVisit) GetParameters() map[string][]string {
	return map[string][]string{
		"PUT": []string{"blog_id:int"},
	}
}

func (this *UserVisit) Put(ctx *eel.Context) {
	req := ctx.Request
	blogId, _ := req.GetInt("blog_id")

	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	blogRepository := blog.NewBlogRepository(bCtx)
	blog := blogRepository.GetBlog(blogId)
	blog.VisitByUser(user)

	ctx.Response.JSON(eel.Map{
	})
}
