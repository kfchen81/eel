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
	UserVisits []*RUserVisit `json:"user_visits"`
	SelfVisit *RSelfVisit `json:"self_visit"`
	CreatedAt string `json:"created_at"`
}

type RUserVisit struct {
	Id int `json:"id"'`
	User *account.RUser `json:"user"`
	Blog *RBlog `json:"blog"`
	CreatedAt string `json:"created_at"`
}

type RSelfVisit struct {
	Id int `json:"id"'`
	CreatedAt string `json:"created_at"`
}

func init() {
}
