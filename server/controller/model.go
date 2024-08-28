package controller

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}
