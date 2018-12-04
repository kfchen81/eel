package blog

import (
	"context"
	"github.com/kfchen81/eel/devapp/business/account"
	m_blog "github.com/kfchen81/eel/devapp/models/blog"
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

	blogIds := make([]int, 0)
	for _, blog := range blogs {
		blogIds = append(blogIds, blog.Id)
	}

	if enableOption, ok := option["with_user"]; ok && enableOption {
		this.fillUser(blogs)
	}

	if enableOption, ok := option["with_user_actions"]; ok && enableOption {
		this.fillUserActions(blogs, blogIds)
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

func (this *FillBlogService) fillUserActions(blogs []*Blog, blogIds []int) {
	o := eel.GetOrmFromContext(this.Ctx)
	
	var models []*m_blog.UserVisit
	err := o.Filter("blog_id__in", blogIds).All(&models)
	
	if err != nil {
		eel.Logger.Error(err)
		return
	}
	
	//构建<id, blog>
	id2blog := make(map[int]*Blog)
	for _, blog := range blogs {
		id2blog[blog.Id] = blog
	}
	
	//获取访问的user信息，构建<userId, user>
	userIds := make([]int, 0)
	for _, model := range models {
		userIds = append(userIds, model.UserId)
	}

	userRepository := account.NewUserRepository(this.Ctx)
	users := userRepository.GetUsers(userIds)
	id2user := make(map[int]*account.User)
	for _, user := range users {
		id2user[user.Id] = user
	}

	//遍历models，设置各个blog的Tags
	for _, model := range models {
		blogId := model.BlogId
		if blog, ok := id2blog[blogId]; ok {
			userVisit := NewUserVisitFromModel(this.Ctx, model)
			if user, ok := id2user[model.UserId]; ok && user != nil {
				userVisit.User = user
			}
			blog.AddUserVisit(userVisit)
		}
	}
}

func init() {
}
