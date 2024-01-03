package ginx

import (
	"github.com/gin-gonic/gin"
	ijwt "goFoundation/webook/internal/web/jwt"
)

func WrapReq[T any](fn func(ctx *gin.Context, req T, claims ijwt.UserClaims) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

type Result struct {
	// 这个叫做业务错误码
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}
