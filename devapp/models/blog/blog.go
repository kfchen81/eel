package blog

import (
	"github.com/kfchen81/eel"
)

type Blog struct {
	eel.Model
	UserId int
	Title string `gorm:"size:1024"`
	Content string `gorm:"size:1024"`
}

func (this *Blog) TableName() string {
	return "blog_blog"
}

func init() {
	eel.RegisterModel(&Blog{})
}