package utils

type BusinessError struct {
	ErrCode string
	ErrMsg string
}

func NewBusinessError(code string, msg string) *BusinessError {
	return &BusinessError{
		code,
		msg,
	}
}

func (this *BusinessError) Error() string {
	return this.ErrCode
}
