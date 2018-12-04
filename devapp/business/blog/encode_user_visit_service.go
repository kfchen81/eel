package blog

import (
	"context"
	"github.com/kfchen81/eel/devapp/business/account"

	"github.com/kfchen81/eel"
)

type EncodeUserVisitService struct {
	eel.ServiceBase
}

func NewEncodeUserVisitService(ctx context.Context) *EncodeUserVisitService {
	service := new(EncodeUserVisitService)
	service.Ctx = ctx
	return service
}

//Encode 对单个实体对象进行编码
func (this *EncodeUserVisitService) Encode(userVisit *UserVisit) *RUserVisit {
	//编码User
	var rUser *account.RUser = nil
	if userVisit.User != nil {
		user := userVisit.User
		
		rUser = &account.RUser{
			Id: user.Id,
			Name: user.Name,
			Avatar: user.Avatar,
			Sex: user.Sex,
		}
	}

	return &RUserVisit{
		Id: userVisit.Id,
		User: rUser,
		CreatedAt: userVisit.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

//EncodeMany 对实体对象进行批量编码
func (this *EncodeUserVisitService) EncodeMany(userVisits []*UserVisit) []*RUserVisit {
	rDatas := make([]*RUserVisit, 0)
	for _, userVisit := range userVisits {
		rDatas = append(rDatas, this.Encode(userVisit))
	}
	
	return rDatas
}

func init() {
}
