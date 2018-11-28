package blog

import (
	"context"
	"github.com/kfchen81/eel/devapp/business"
	m_blog "github.com/kfchen81/eel/devapp/models/blog"

	"github.com/kfchen81/eel"
	"github.com/kfchen81/gorm"
)

type BlogRepository struct {
	eel.RepositoryBase
}

func NewBlogRepository(ctx context.Context) *BlogRepository {
	repository := new(BlogRepository)
	repository.Ctx = ctx
	return repository
}

//GetBlogs 获得Blog对象集合
func (this *BlogRepository) GetBlogsForUser(user business.IUser, page *eel.PageInfo) ([]*Blog, eel.INextPageInfo) {
	o := eel.GetOrmFromContext(this.Ctx)
	qs := o.QueryTable(&m_blog.Blog{})

	var models []*m_blog.Blog
	qs = qs.Filter("user_id", user.GetId()).Filter("is_deleted", false).OrderBy("-id")
	paginateResult, err := eel.Paginate(qs, page, &models)

	if err != nil {
		eel.Logger.Error(err)
		return nil, paginateResult
	}

	blogs := make([]*Blog, 0)
	for _, model := range models {
		blogs = append(blogs, NewBlogFromModel(this.Ctx, model))
	}
	return blogs, paginateResult
}

//GetBlog 根据id获得Blog对象
func (this *BlogRepository) GetBlog(id int) *Blog {
	var mBlog m_blog.Blog

	err := eel.GetOrmFromContext(this.Ctx).Filter("id", id).One(&mBlog)

	if err != nil {
		eel.Logger.Error(err)
		return nil
	}

	blog := NewBlogFromModel(this.Ctx, &mBlog)
	return blog
}

//DeleteBlog 根据id删除Blog对象
func (this *BlogRepository) DeleteBlog(id int) bool {
	db := eel.GetOrmFromContext(this.Ctx).QueryTable(&m_blog.Blog{})

	//_, err := qs.Filter("id", id).Delete()
	db.Filter("id", id).Update(gorm.Params{
		"is_deleted": true,
	})

	return true
}

func init() {
}
