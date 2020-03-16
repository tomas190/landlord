package controller

const (
	httpCode = 200
	SuccCode = 0
	ErrCode  = -1
)

type ApiResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewResp(code int, msg string, data interface{}) ApiResp {
	return ApiResp{Code: code, Msg: msg, Data: data}
}
