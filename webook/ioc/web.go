package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"goFoundation/webook/internal/web"
	"goFoundation/webook/internal/web/middleware"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRouter(server)
	return server
}

// 中间件校验token
func InitMiddlewares(rdb redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		corsHdl(),
		middleware.NewUserJwtService().
			LoginPath("/users/signup").
			LoginPath("/users/login").
			LoginPath("/users/login_sms/code/send").
			LoginPath("/users/login_sms").
			Build(),
		//ratelimit.NewBuilder(rdb, time.Second, 100).Build(),
	}
}

// 跨域
func corsHdl() gin.HandlerFunc {
	return cors.New(cors.Config{
		//AllowOrigins:     []string{"http://localhost:3000"},         //前端地址
		AllowCredentials: true,                                      //是否携带cookie
		AllowHeaders:     []string{"Content-Type", "Authorization"}, //字段类型json类型
		ExposeHeaders:    []string{"x-jwt-token"},                   //用jwt方法发送token必须跨域时配置这个 否则前端无法获取token
		//AllowMethods:     []string{"POST"},                          //什么方法post,get
		AllowOriginFunc: func(origin string) bool {
			//开发环境
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.Contains(origin, "公司域名")
		},
		MaxAge: 12 * time.Hour,
	})
}
