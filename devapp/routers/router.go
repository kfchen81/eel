package routers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/kfchen81/eel/log"
	"github.com/kfchen81/eel/router"
	"github.com/kfchen81/eel/devapp/rest/blog"
)

func init() {
	router.RegisterResource(&blog.Blog{})
	log.Info("in router...")
}