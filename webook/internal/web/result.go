package web

import "github.com/gin-gonic/gin"

type Result struct {
	// 这个叫做业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type handler interface {
	RegisterRoutes(s *gin.Engine)
}
