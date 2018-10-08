package console

import (
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/utils"
	"html/template"
	"bytes"
	"github.com/kfchen81/eel/router"
	"github.com/kfchen81/eel/config"
)

type Console struct {
	eel.RestResource
}

func (this *Console) Resource() string {
	return "console.console"
}

func (this *Console) GetParameters() map[string][]string {
	return map[string][]string{
		"GET":    []string{},
	}
}

func (this *Console) Get(ctx *eel.Context) {
	path := utils.SearchFileInGoPath("static/service_console.html")
	t, _ := template.ParseFiles(path)
	var bufferWriter bytes.Buffer
	resources := router.Resources()
	
	serviceName := config.ServiceConfig.String("appname")
	t.Execute(&bufferWriter, map[string]interface{}{
		"Resources": resources,
		"Name": serviceName,
	})
	
	ctx.Response.Content(bufferWriter.Bytes(), "text/html; charset=utf-8")
}
