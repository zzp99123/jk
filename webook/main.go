package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"goFoundation/webook/config"
	"goFoundation/webook/internal/repository"
	"goFoundation/webook/internal/repository/dao"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/internal/web"
	"goFoundation/webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	d := initDb()
	server := initWebServer()
	rdb := initRedis()
	u := initUser(d, rdb)
	u.RegisterRouter(server)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})
	server.Run(":8081")
}
func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
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
	}))
	//cookie 不安全
	//store := cookie.NewStore([]byte("secret"))
	//memsotre 单实例部署 很少用
	//参数放的是密钥
	store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	//redis 多实例部署 无脑用
	//第一个参数空闲连接数量 第二个参数 tcp连接 第三个连接信息 第四个密码 第五六就是key
	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
	//	[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	//if err != nil {
	//	panic(err)
	//}
	server.Use(sessions.Sessions("ssid", store))
	//登录检验看我有么有这个id
	server.Use(middleware.NewUserService().
		//链式调用
		LoginPath("/users/signup").
		LoginPath("/users/login").Build())
	//jwt
	//server.Use(middleware.NewUserJwtService().
	//	LoginPath("/users/signup").
	//	LoginPath("/users/login").Build())
	return server
}
func initUser(db *gorm.DB, rdb redis.Cmdable) *web.UserHandler {
	d := dao.NewUserService(db)
	r := repository.NewUserService(d)
	s := service.NewUserService(r)
	u := web.NewUserHandle(s)
	return u
}
func initDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

func initRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})
	return redisClient
}
