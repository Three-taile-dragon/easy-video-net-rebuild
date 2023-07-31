package common

import "net/http"

// 定义消息返回体

type BusinessCode int
type Result struct {
	Code    BusinessCode `json:"code"`
	Message string       `json:"message"`
	Data    any          `json:"data"`
}

func (r *Result) Success(data any) *Result {
	r.Code = http.StatusOK
	r.Message = "success"
	r.Data = data
	return r
}

func (r *Result) Fail(code BusinessCode, msg string) *Result {
	r.Code = code
	r.Message = msg
	return r
}
