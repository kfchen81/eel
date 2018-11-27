package blog

import (
	"fmt"
	"context"
	"github.com/kfchen81/eel/devapp/business"
	"github.com/kfchen81/eel/devapp/business/account"
	m_blog "github.com/kfchen81/eel/devapp/models/blog"
	"time"

	"github.com/kfchen81/eel"
	"github.com/kfchen81/gorm"
)

type Blog struct {
	eel.EntityBase
	Id int
	UserId int
	Title string
	Content string
	IsDeleted bool
	CreatedAt time.Time

	User *account.User
}

//Update 更新对象
func (this *Blog) Update(title string, content string) {
	var mBlog m_blog.Blog
	o := eel.GetOrmFromContext(this.Ctx)

	o.QueryTable(&mBlog).Filter("id", this.Id).Update(gorm.Params{
		"title": title,
		"content": content,
	})
}

//工厂方法
func NewBlog(ctx context.Context, user business.IUser, title string, content string) *Blog {
	o := eel.GetOrmFromContext(ctx)
	mBlog := m_blog.Blog{}
	mBlog.UserId = user.GetId()
	mBlog.Title = title
	mBlog.Content = content
	result := o.Insert(&mBlog)
	if result.Error != nil {
		eel.Logger.Error(result.Error)
		panic(eel.NewBusinessError("blog:create_fail", fmt.Sprintf("创建失败")))
	}

	return NewBlogFromModel(ctx, &mBlog)
}

//根据model构建对象
func NewBlogFromModel(ctx context.Context, mBlog *m_blog.Blog) *Blog {
	blog := new(Blog)
	blog.Ctx = ctx
	blog.Model = mBlog
	blog.Id = mBlog.Id
	blog.UserId = mBlog.UserId
	blog.Title = mBlog.Title
	blog.Content = mBlog.Content
	blog.IsDeleted = mBlog.IsDeleted
	blog.CreatedAt = mBlog.CreatedAt
	return blog
}

func init() {
}
