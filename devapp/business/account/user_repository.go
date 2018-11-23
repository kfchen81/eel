package account

import (
	"context"
	"encoding/json"
	
	"fmt"
	"github.com/kfchen81/eel"
)

type UserRepository struct {
	eel.ServiceBase
}

func NewUserRepository(ctx context.Context) *UserRepository {
	service := new(UserRepository)
	service.Ctx = ctx
	return service
}

func (this *UserRepository) makeUsers(userDatas []interface{}) []*User {
	users := make([]*User, 0)
	for _, userData := range userDatas {
		userJson := userData.(map[string]interface{})
		id, _ := userJson["id"].(json.Number).Int64()
		user := NewUserFromOnlyId(this.Ctx, int(id))
		user.Name = userJson["name"].(string)
		user.Avatar = userJson["avatar"].(string)
		user.Cover = userJson["cover"].(string)
		user.Birthday = userJson["birthday"].(string)
		user.Region = userJson["region"].(string)
		user.Slogan = userJson["slogan"].(string)
		user.Sex = userJson["sex"].(string)
		user.Code = userJson["code"].(string)
		user.Phone = userJson["phone"].(string)
		user.IsRegisterEasemob = userJson["is_register_easemob"].(bool)
		user.Age, _ = userJson["age"].(json.Number).Int64()
		user.Latitude, _ = userJson["latitude"].(json.Number).Float64()
		user.Longitude, _ = userJson["longitude"].(json.Number).Float64()
		
		if roles, ok := userJson["roles"]; ok {
			user.Roles = roles.([]interface{})
		}
		users = append(users, user)
	}
	
	return users
}

func (this *UserRepository) GetUsers(ids []int) []*User {
	options := make(map[string]interface{})
	options["with_role_info"] = true
	resp, err := eel.NewResource(this.Ctx).Get("gskep", "account.users", eel.Map{
		"ids": eel.ToJsonString(ids),
		"with_options": eel.ToJsonString(options),
	})
	
	if err != nil {
		eel.Logger.Error(err)
		return nil
	}
	
	respData := resp.Data()
	userDatas := respData.Get("users")
	fmt.Println(userDatas)
	return this.makeUsers(userDatas.MustArray())
}

func (this *UserRepository) GetUsersWithOptions(ids []int, options map[string]interface{}) []*User {
	resp, err := eel.NewResource(this.Ctx).Get("gskep", "account.users", eel.Map{
		"ids": eel.ToJsonString(ids),
		"with_options": eel.ToJsonString(options),
	})
	
	if err != nil {
		eel.Logger.Error(err)
		return nil
	}
	
	respData := resp.Data()
	userDatas := respData.Get("users")
	return this.makeUsers(userDatas.MustArray())
}

func (this *UserRepository) GetUsersByCodes(codes []string) []*User {
	resp, err := eel.NewResource(this.Ctx).Get("gskep", "account.users", eel.Map{
		"codes": eel.ToJsonString(codes),
	})
	
	if err != nil {
		eel.Logger.Error(err)
		return nil
	}
	
	respData := resp.Data()
	userDatas := respData.Get("users")
	return this.makeUsers(userDatas.MustArray())
}

func (this *UserRepository) GetUsersByUnionids(unionids []string) []*User {
	resp, err := eel.NewResource(this.Ctx).Get("gskep", "account.users", eel.Map{
		"unionids": eel.ToJsonString(unionids),
	})
	
	if err != nil {
		eel.Logger.Error(err)
		return nil
	}
	
	respData := resp.Data()
	userDatas := respData.Get("users")
	return this.makeUsers(userDatas.MustArray())
}

func init() {
}
