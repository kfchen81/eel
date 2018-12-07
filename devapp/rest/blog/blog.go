package blog

import (
	"github.com/kfchen81/eel/devapp/business/account"
	b_blog "github.com/kfchen81/eel/devapp/business/blog"

	"github.com/kfchen81/eel"
)

type Blog struct {
	eel.RestResource
}

func (this *Blog) Resource() string {
	return "blog.blog"
}

func (this *Blog) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{"id:int"},
		"PUT": []string{"title", "content", },
		"POST": []string{"id:int", "title", "content", },
		"DELETE": []string{"id:int"},
	}
}

func (this *Blog) Get(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	blogRepository := b_blog.NewBlogRepository(bCtx)
	blog := blogRepository.GetBlog(id)

	fillService := b_blog.NewFillBlogService(bCtx)
	fillService.Fill([]*b_blog.Blog{ blog }, eel.FillOption{
		"with_user": true,
		"with_user_actions": true,
	})

	encodeService := b_blog.NewEncodeBlogService(bCtx)
	respData := encodeService.Encode(blog)

	ctx.Response.JSON(respData)
}

func (this *Blog) Put(ctx *eel.Context) {
	req := ctx.Request
	title := req.GetString("title")
	content := req.GetString("content")

	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	blog := b_blog.NewBlog(bCtx, user, title, content)

	ctx.Response.JSON(eel.Map{
		"id": blog.Id,
	})
}

func (this *Blog) Post(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")
	title := req.GetString("title")
	content := req.GetString("content")

	bCtx := ctx.GetBusinessContext()
	blogRepository := b_blog.NewBlogRepository(bCtx)
	blog := blogRepository.GetBlog(id)

	blog.Update(title, content)

	ctx.Response.JSON(eel.Map{})
}

func (this *Blog) Delete(ctx *eel.Context) {
	req := ctx.Request
	id, _ := req.GetInt("id")

	bCtx := ctx.GetBusinessContext()
	blogRepository := b_blog.NewBlogRepository(bCtx)
	blogRepository.DeleteBlog(id)

	ctx.Response.JSON(eel.Map{})
}
