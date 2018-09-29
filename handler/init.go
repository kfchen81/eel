package handler

import "fmt"

//returnValidateParameterFailResponse 返回参数校验错误的response
func returnValidateParameterFailResponse(ctx *Context, parameter string, paramType string, innerErrMsgs ...string) {
	innerErrMsg := ""
	if len(innerErrMsgs) > 0 {
		innerErrMsg = innerErrMsgs[0]
	}
	
	ctx.Response.ErrorWithCode(
		500,
		"rest:missing_argument",
		fmt.Sprintf("missing or invalid argument: %s(%s)", parameter, paramType),
		innerErrMsg,
	)
}

func init() {

}