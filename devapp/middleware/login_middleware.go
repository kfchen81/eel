package middleware

import (
	"github.com/kfchen81/eel"
	"github.com/kfchen81/eel/log"
)

type LoginMiddleware struct {
	eel.Middleware
}

func (this *LoginMiddleware) ProcessRequest(ctx *eel.Context) {
	log.Info("i am in login middleware process request")
}

