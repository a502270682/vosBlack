package error_code

type ReplyError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Error(code int, msg string) *ReplyError {
	return &ReplyError{
		Code:    code,
		Message: msg,
	}
}
