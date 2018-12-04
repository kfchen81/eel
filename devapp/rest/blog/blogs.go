package blog

import (
	"github.com/kfchen81/eel/devapp/business/blog"
	"github.com/kfchen81/eel/devapp/business/account"

	"github.com/kfchen81/eel"
)

type Blogs struct {
	eel.RestResource
}

func (this *Blogs) Resource() string {
	return "blog.blogs"
}

func (this *Blogs) GetParameters() map[string][]string {
	return map[string][]string{
		"GET": []string{},
	}
}

func (this *Blogs) Get(ctx *eel.Context) {
	bCtx := ctx.GetBusinessContext()
	user := account.GetUserFromContext(bCtx)
	page := eel.ExtractPageInfoFromRequest(ctx)
	repository := blog.NewBlogRepository(bCtx)
	blogs, nextPageInfo := repository.GetBlogsForUser(user, page)

	fillService := blog.NewFillBlogService(bCtx)
	fillService.Fill(blogs, eel.FillOption{
		"with_user": true,
		
	})

	encodeService := blog.NewEncodeBlogService(bCtx)
	rows := make([]*blog.RBlog, 0)
	for _, blog := range blogs {
		rows = append(rows, encodeService.Encode(blog))
	}

	ctx.Response.JSON(eel.Map{
		"blogs": rows,
		"pageinfo": nextPageInfo.ToMap(),
	})
}
