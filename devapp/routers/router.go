package routers

import (
	"github.com/kfchen81/eel/devapp/rest/blog"
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/handler/rest/console"
)

func init() {
	eel.RegisterResource(&console.Console{})
	eel.RegisterResource(&blog.Blog{})
	eel.RegisterResource(&blog.Blogs{})
	//eel.RegisterResource(&blog.BlogTags{})
}