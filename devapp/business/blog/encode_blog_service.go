package blog

import (
	"context"
	"github.com/kfchen81/eel/devapp/business/account"

	"github.com/kfchen81/eel"
)

type EncodeBlogService struct {
	eel.ServiceBase
}

func NewEncodeBlogService(ctx context.Context) *EncodeBlogService {
	service := new(EncodeBlogService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeBlogService) Encode(blog *Blog) *RBlog {
	//编码User
	var rUser *account.RUser = nil
	if blog.User != nil {
		user := blog.User
		
		rUser = &account.RUser{
			Id:     user.Id,
			Name:   user.Name,
			Avatar: user.Avatar,
			Sex:    user.Sex,
		}
	}

	//编码UserVisit
	var rUserVisits []*RUserVisit
	var rSelfVisit *RSelfVisit
	if blog.HasUserVisits() {
		encodeUserVisitService := NewEncodeUserVisitService(this.Ctx)
		rUserVisits = encodeUserVisitService.EncodeMany(blog.UserVisits)

		selfVisit := blog.GetSelfVisit()
		if selfVisit != nil {
			rSelfVisit = &RSelfVisit{
				Id: selfVisit.Id,
				CreatedAt: selfVisit.CreatedAt.Format("2006-01-02 15:04:05"),
			}
		}
	}

	return &RBlog{
		Id: blog.Id,
		User: rUser,
		UserVisits: rUserVisits,
		SelfVisit: rSelfVisit,
		Title: blog.Title,
		Content: blog.Content,
		IsDeleted: blog.IsDeleted,
		CreatedAt: blog.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeBlogService) EncodeMany(blogs []*Blog) []*RBlog {
	rDatas := make([]*RBlog, 0)
	for _, blog := range blogs {
		rDatas = append(rDatas, this.Encode(blog))
	}
	
	return rDatas
}

func init() {
}
