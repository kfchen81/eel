package account

type RUser struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Cover  string `json:"cover"`
	Sex    string `json:"sex"`
	Code   string `json:"code"`
	Phone   string `json:"phone"`
	Birthday string `json:"birthday"`
	Region string `json:"region"`
	Slogan string `json:"slogan"`
	Age int64 `json:"age"`
	Roles []string `json:"roles"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	IsFollowee bool `json:"is_followee"`
	IsRegisterEasemob bool `json:"is_register_easemob"`
}
