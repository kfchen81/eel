package blog

import (
	"github.com/kfchen81/eel"
)

type Blogs struct {
	eel.RestResource
}

func (this *Blogs) Resource() string {
	return "blog.blogs"
}