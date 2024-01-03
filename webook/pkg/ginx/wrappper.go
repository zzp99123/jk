package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/prometheus/client_golang/prometheus"
	"goFoundation/webook/pkg/logger"
	"net/http"
	"strconv"
)

// 这个东西，放到你们的 ginx 插件库里面去
// 技术含量不是很高，但是绝对有技巧

// L 使用包变量
var L logger.LoggerV1
var vec *prometheus.CounterVec

func NewCounterVec(opt prometheus.CounterOpts) {
	vector := prometheus.NewCounterVec(opt, []string{"type"})
	prometheus.MustRegister(vector)
}

func WrapToken[C jwt.Claims](fn func(ctx *gin.Context, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val, ok := ctx.Get("claims")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, claims)
		fmt.Println(res, "7777777777777777777777777777777777777777")
		if err != nil {
			L.Error("处理业务逻辑出错",
				logger.String("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		vec.WithLabelValues(strconv.Itoa(res.Code)).Inc()
		ctx.JSON(http.StatusOK, res)
	}
}
func WrapBodyAndToken[Req any, C jwt.Claims](fn func(ctx *gin.Context, r Req, uc C) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var r Req
		val, ok := ctx.Get("claims")
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		claims, ok := val.(C)
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		res, err := fn(ctx, r, claims)
		if err != nil {
			L.Error("处理业务逻辑出错",
				logger.String("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}
func WrapBodyV1[Req any](fn func(ctx *gin.Context, r Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var r Req
		res, err := fn(ctx, r)
		// 下半段的业务逻辑从哪里来？
		// 我的业务逻辑有可能要操作 ctx
		// 你要读取 HTTP HEADER
		if err != nil {
			L.Error("处理业务逻辑出错",
				logger.String("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}

func WrapBody[Req any](l logger.LoggerV1, fn func(ctx *gin.Context, r Req) (Result, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var r Req
		// 下半段的业务逻辑从哪里来？
		// 我的业务逻辑有可能要操作 ctx
		// 你要读取 HTTP HEADER
		res, err := fn(ctx, r)
		if err != nil {
			// 开始处理 error，其实就是记录一下日志
			l.Error("处理业务逻辑出错",
				logger.String("path", ctx.Request.URL.Path),
				// 命中的路由
				logger.String("route", ctx.FullPath()),
				logger.Error(err))
		}
		ctx.JSON(http.StatusOK, res)
	}
}
