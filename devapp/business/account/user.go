package account

import (
	"context"

	"time"

	//	"github.com/astaxie/eel"

	"github.com/bitly/go-simplejson"
	
	"github.com/kfchen81/eel"
)


type User struct {
	eel.EntityBase
	Id                int
	PlatformId        int
	Name              string
	Avatar            string
	Cover             string
	Sex               string
	Phone             string
	Birthday          string
	Region            string
	Slogan            string
	Longitude         float64
	Latitude          float64
	Distance          string
	Age               int64
	Code              string
	RawData           *simplejson.Json
	LastActiveTime    string
	DisplayLiveness   string
	IsRegisterEasemob bool
	Roles             []interface{}
	CreatedAt         time.Time
}

func NewUserFromOnlyId(ctx context.Context, id int) *User {
	user := new(User)
	user.Ctx = ctx
	user.Model = nil
	user.Id = id
	return user
}

func (this *User) GetId() int {
	return this.Id
}

func GetUserFromContext(ctx context.Context) *User {
	user := ctx.Value("user").(*User)
	return user
}

func init() {
}
