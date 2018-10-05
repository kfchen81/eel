package handler

import (
	"fmt"
	"bytes"
	"runtime"
	"github.com/kfchen81/eel/util"
	"github.com/kfchen81/eel/log"
)

func RecoverPanic(ctx *Context) {
	if err := recover(); err != nil {
		//rollback commit
		//o := ctx.Input.Data()["sessionOrm"]
		//o.(orm.Ormer).Rollback()
		log.Warn("[ORM] rollback transaction 2")
		
		//finish span
		//span := ctx.Input.GetData("span")
		//if span != nil {
		//	beego.Info("[Tracing] finish span in recoverPanic")
		//	span.(opentracing.Span).Finish()
		//}
		
		errMsg := ""
		if be, ok := err.(*util.BusinessError); ok {
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
		log.Error(buffer.String())
		
		var resp Map
		if be, ok := err.(*util.BusinessError); ok {
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
	}
}
