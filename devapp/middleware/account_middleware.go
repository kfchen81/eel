package middleware

import (
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/log"
)

type AccountMiddleware struct {
	eel.Middleware
}

func (this *AccountMiddleware) ProcessRequest(ctx *eel.Context) {
	ctx.Set("account", "陈浩宇")
	log.Info("i am in account middleware process request")
}

func (this *AccountMiddleware) ProcessResponse(ctx *eel.Context) {
	log.Info("i am in account middleware process response")
}

