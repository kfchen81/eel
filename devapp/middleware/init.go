package middleware

import "github.com/kfchen81/eel"

func init() {
	eel.RegisterMiddleware(&AccountMiddleware{})
	eel.RegisterMiddleware(&LoginMiddleware{})
}
