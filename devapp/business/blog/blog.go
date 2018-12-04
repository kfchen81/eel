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
	UserVisits []*UserVisit
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

//VisitByUser user访问Blog
func (this *Blog) VisitByUser(user business.IUser) {
	o := eel.GetOrmFromContext(this.Ctx)

	isExist := o.QueryTable(&m_blog.UserVisit{}).Filter(gorm.Map{
		"blog_id": this.Id,
		"user_id": user.GetId(),
	}).Exist()

	if !isExist {
		model := &m_blog.UserVisit{}
		model.BlogId = this.Id
		model.UserId = user.GetId()
		err := o.Insert(&model)

		if err != nil {
			eel.Logger.Error(err)
			panic(eel.NewBusinessError("blog:user_visit_fail", "创建用户访问记录失败"))
		}
	}
}

func (this *Blog) AddUserVisit(userVisit *UserVisit) {
	this.UserVisits = append(this.UserVisits, userVisit)
}

func (this *Blog) HasUserVisits() bool {
	return len(this.UserVisits) > 0
}

func (this *Blog) IsVisitedByUser(user business.IUser) bool {
	if len(this.UserVisits) == 0 {
		//没有填充UserVisits时，通过数据库判断
		o := eel.GetOrmFromContext(this.Ctx)

		isExist := o.QueryTable(&m_blog.UserVisit{}).Filter(gorm.Map{
			"blog_id": this.Id,
			"user_id": user.GetId(),
		}).Exist()

		return isExist
	} else {
		targetId := this.Id
		for _, userVisit := range this.UserVisits {
			if userVisit.BlogId == targetId {
				return true
			}
		}

		return false
	}
	return len(this.UserVisits) > 0
}

func (this *Blog) GetSelfVisit() *UserVisit {
	if len(this.UserVisits) == 0 {
		return nil
	} else {
		user := account.GetUserFromContext(this.Ctx)
		targetUserId := user.Id
		for _, userVisit := range this.UserVisits {
			if userVisit.UserId == targetUserId {
				return userVisit
			}
		}

		return nil
	}
}

//工厂方法
func NewBlog(ctx context.Context, user business.IUser, title string, content string) *Blog {
	o := eel.GetOrmFromContext(ctx)
	mBlog := m_blog.Blog{}
	mBlog.UserId = user.GetId()
	mBlog.Title = title
	mBlog.Content = content
	err := o.Insert(&mBlog)
	if err != nil {
		eel.Logger.Error(err)
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
	blog.UserVisits = make([]*UserVisit, 0)

	return blog
}

func init() {
}
