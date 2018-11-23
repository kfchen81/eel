package handler

type RestResponse struct {
	Code        int32                  `json:"code"`
	Data        interface{} `json:"data"`
	ErrCode string `json:"errCode"`
	ErrMsg      string                 `json:"errMsg"`
	InnerErrMsg string                 `json:"innerErrMsg"`
}


func MakeResponse(data interface{}) *RestResponse {
	return &RestResponse{
		200,
		data,
		"",
		"",
		"",
	}
}

func MakeErrorResponse(code int32, errCode string, errMsg string, innerErrMsgs ...string) *RestResponse {
	innerErrMsg := ""
	if len(innerErrMsgs) > 0 {
		innerErrMsg = innerErrMsgs[0]
	}
	
	return &RestResponse{
		code,
		nil,
		errCode,
		errMsg,
		innerErrMsg,
	}
}