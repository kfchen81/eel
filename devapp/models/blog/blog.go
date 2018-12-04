package blog

import (
	"github.com/kfchen81/eel"
)

//Blog Model
type Blog struct {
	eel.DeletableModel
	UserId int `gorm:"index"`//foreign key for user
	Title string `gorm:"size:1024"`
	Content string `gorm:"size:1024"`
}
func (self *Blog) TableName() string {
	return "blog_blog"
}


//UserVisit Model
type UserVisit struct {
	eel.Model
	BlogId int `gorm:"index"`
	UserId int `gorm:"index"`
}
func (self *UserVisit) TableName() string {
	return "blog_user_visit"
}


func init() {
	eel.RegisterModel(new(Blog))
	eel.RegisterModel(new(UserVisit))
}
