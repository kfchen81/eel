package blog

import (
	"context"
	m_blog "github.com/kfchen81/eel/devapp/models/blog"
	"github.com/kfchen81/eel"
	"time"
	"github.com/kfchen81/eel/devapp/business"
	"github.com/pkg/errors"
)

type Blog struct {
	eel.EntityBase
	Id int
	Title string
	Content string
	CreatedAt time.Time
}

func NewBlog(ctx context.Context, user business.IUser, title string, content string) *Blog {
	model := m_blog.Blog{
		UserId: user.GetId(),
		Title: title,
		Content: content,
	}
	
	orm := eel.GetOrmFromContext(ctx)
	orm.Create(&model)
	
	panic(errors.New("lala haha hehe huhu"))
	
	return NewBlogFromModel(ctx, &model)
}

//根据model构建对象
func NewBlogFromModel(ctx context.Context, mBlog *m_blog.Blog) *Blog {
	blog := new(Blog)
	blog.Ctx = ctx
	blog.Model = mBlog
	blog.Id = mBlog.ID
	blog.Title = mBlog.Title
	blog.Content = mBlog.Content
	blog.CreatedAt = mBlog.CreatedAt
	return blog
}