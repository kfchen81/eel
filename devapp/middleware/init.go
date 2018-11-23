package middleware

import "github.com/kfchen81/eel"

func init() {
	eel.RegisterMiddleware(&eel.JWTMiddleware{})
	eel.RegisterMiddleware(&LoginMiddleware{})
}
