package blog

import (
	"context"
	"github.com/kfchen81/eel/devapp/business/account"
	m_blog "github.com/kfchen81/eel/devapp/models/blog"
	"time"

	"github.com/kfchen81/eel"
)

type UserVisit struct {
	eel.EntityBase
	Id int
	BlogId int
	UserId int
	CreatedAt time.Time

	User *account.User
	Blog *Blog
}

//根据model构建对象
func NewUserVisitFromModel(ctx context.Context, model *m_blog.UserVisit) *UserVisit {
	instance := new(UserVisit)
	instance.Ctx = ctx
	instance.Model = model
	instance.Id = model.Id
	instance.BlogId = model.BlogId
	instance.UserId = model.UserId

	return instance
}

func init() {
}
