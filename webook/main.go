package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//d := initDb()
	//server := initWebServer()
	//rdb := initRedis()
	//u := initUser(d, rdb)
	//u.RegisterRouter(server)
	server := InitWebServer()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8080")
}

//func initWebServer() *gin.Engine {
//	//cookie 不安全
//	//store := cookie.NewStore([]byte("secret"))
//	//memsotre 单实例部署 很少用
//	//参数放的是密钥
//	//store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
//	//redis 多实例部署 无脑用
//	//第一个参数空闲连接数量 第二个参数 tcp连接 第三个连接信息 第四个密码 第五六就是key
//	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
//	//	[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	//server.Use(sessions.Sessions("ssid", store))
//	//登录检验看我有么有这个id
//	//server.Use(middleware.NewUserService().
//	//	//链式调用
//	//	LoginPath("/users/signup").
//	//	LoginPath("/users/login").Build())
//	return server
//}
