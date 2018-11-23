package handler

import (
	"fmt"
	"bytes"
	"runtime"
	"github.com/kfchen81/eel/utils"
	"github.com/kfchen81/eel/log"
	"github.com/kfchen81/gorm"
)

func RecoverPanic(ctx *Context) {
	log.Logger.Debug("[router] in RecoverPanic...")
	if err := recover(); err != nil {
		orm := ctx.Get("orm")
		if orm != nil {
			log.Logger.Warn("[ORM] rollback transaction")
			orm.(*gorm.DB).Rollback()
		}
		
		errMsg := ""
		if be, ok := err.(*utils.BusinessError); ok {
			errMsg = fmt.Sprintf("%s:%s", be.ErrCode, be.ErrMsg)
		
		} else {
			errMsg = fmt.Sprintf("%s", err)
		}
		
		var buffer bytes.Buffer
		buffer.WriteString(fmt.Sprintf("[Unprocessed_Exception] %s\n", errMsg))
		buffer.WriteString(fmt.Sprintf("Request URL: %s\n", ctx.Request.URL()))
		for i := 1; ; i++ {
			_, file, line, ok := runtime.Caller(i)
			if !ok {
				break
			}
			buffer.WriteString(fmt.Sprintf("%s:%d\n", file, line))
		}
		log.Logger.Error(buffer.String())
		
		var resp Map
		if be, ok := err.(*utils.BusinessError); ok {
			resp = Map{
				"code": 500,
				"data": nil,
				"errCode": be.ErrCode,
				"errMsg": be.ErrMsg,
				"innerErrMsg": "",
			}
		} else {
			resp = Map{
				"code": 531,
				"data": nil,
				"errCode": "system:exception",
				"errMsg": fmt.Sprintf("%s", err),
				"innerErrMsg": "",
			}
		}
		ctx.Response.JSON(resp)
	} else {
		orm := ctx.Get("orm")
		if orm != nil {
			log.Logger.Debug("[ORM] commit transaction")
			orm.(*gorm.DB).Commit()
		}
	}
}
