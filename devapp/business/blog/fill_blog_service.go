package blog

import (
	"context"
	"github.com/kfchen81/eel/devapp/business/account"
	"github.com/kfchen81/eel"
)

type FillBlogService struct {
	eel.ServiceBase
}

func NewFillBlogService(ctx context.Context) *FillBlogService {
	service := new(FillBlogService)
	service.Ctx = ctx
	return service
}

func (this *FillBlogService) Fill(blogs []*Blog, option eel.FillOption) {
	if len(blogs) == 0 {
		return
	}

	if _, ok := option["with_user"]; ok {
		this.fillUser(blogs)
	}
}

func (this *FillBlogService) fillUser(blogs []*Blog) {
	userIds := make([]int, 0)
	for _, blog := range blogs {
		userIds = append(userIds, blog.UserId)
	}

	userRepository := account.NewUserRepository(this.Ctx)
	users := userRepository.GetUsers(userIds)
	id2user := make(map[int]*account.User)
	for _, user := range users {
		id2user[user.Id] = user
	}

	for _, blog := range blogs {
		user := id2user[blog.UserId]
		blog.User = user
	}
}

func init() {
}
