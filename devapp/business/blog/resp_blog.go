package blog

import (
	"github.com/kfchen81/eel/devapp/business/account"
)

type RBlog struct {
	Id        int    `json:"id"'`
	Title     string `json:"title"`
	Content   string `json:"content"`
	IsDeleted bool   `json:"is_deleted"`
	User    *account.RUser `json:"author"`
	
	CreatedAt string `json:"created_at"`
}

func init() {
}
