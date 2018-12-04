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
	Counter int
}
func (self *Blog) TableName() string {
	return "blog_blog"
}

func init() {
	eel.RegisterModel(new(Blog))
}
