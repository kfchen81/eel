package routers

import (
	"github.com/kfchen81/eel/devapp/rest/blog"
	"github.com/kfchen81/eel"
)

func init() {
	eel.RegisterResource(&blog.Blog{})
}