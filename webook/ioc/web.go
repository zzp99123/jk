package ioc

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"goFoundation/webook/internal/web"
	ijwt "goFoundation/webook/internal/web/jwt"
	"goFoundation/webook/internal/web/middleware"
	"goFoundation/webook/pkg/ginx"
	"goFoundation/webook/pkg/ginx/middleware/metric"
	logger2 "goFoundation/webook/pkg/logger"
	"strings"
	"time"
)

func InitWebServer(funcs []gin.HandlerFunc, userHdl *web.UserHandler, owh *web.OAuth2WechatHandler, articleHdl *web.ArticleHandler) *gin.Engine {
	server := gin.Default()
	server.Use(funcs...)
	userHdl.RegisterRouter(server)
	owh.RegisterRouters(server)
	articleHdl.RegisterRoutes(server)
	(&web.ObservabilityWeb{}).RegisterRoutes(server)
	return server
}

// 中间件校验token
func InitMiddlewares(rdb redis.Cmdable, jwtHdl ijwt.Handler, l logger2.LoggerV1) []gin.HandlerFunc {
	//打印http接口日志
	//bd := logger.NewBuilder(func(ctx context.Context, al *logger.AccessLog) {
	//	l.Debug("HTTP请求", logger2.Field{Key: "al", Value: al})
	//}).AllowReqBody(true).AllowRespBody()
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//	ok := viper.GetBool("web.logreq")
	//	bd.AllowReqBody(ok)
	//})
	//HTTP 的业务错误码
	ginx.NewCounterVec(prometheus.CounterOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "http_biz_code",
		Help:      "HTTP 的业务错误码",
	})
	//统计 GIN 的 HTTP 接口
	return []gin.HandlerFunc{
		corsHdl(),
		//bd.Build(),
		(&metric.PrometheusMetric{
			Namespace:  "geekbang_daming",
			Subsystem:  "webook",
			Name:       "gin_http",
			Help:       "统计 GIN 的 HTTP 接口",
			InstanceId: "my-instance-1",
		}).Build(),
		otelgin.Middleware("webook"),
		middleware.NewUserJwtService(jwtHdl).
			LoginPath("/users/signup").
			LoginPath("/users/login").
			LoginPath("/users/login_sms/code/send").
			LoginPath("/users/login_sms").
			LoginPath("/oauth2/wechat/authurl").
			LoginPath("/oauth2/wechat/callback").
			LoginPath("/users/refresh_token").
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
		ExposeHeaders:    []string{"x-jwt-token", "x-ref-token"},    //用jwt方法发送token必须跨域时配置这个 否则前端无法获取token
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

//func sessionHandlerFunc() gin.HandlerFunc {
//	//cookie 不安全
//	//store := cookie.NewStore([]byte("secret"))
//	//memsotre 单实例部署 很少用
//	//参数放的是密钥
//	store := memstore.NewStore([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
//	//redis 多实例部署 无脑用
//	//第一个参数空闲连接数量 第二个参数 tcp连接 第三个连接信息 第四个密码 第五六就是key
//	//store, err := redis.NewStore(16, "tcp", "localhost:6379", "",
//	//	[]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
//	//if err != nil {
//	//	panic(err)
//	//}
//	return sessions.Sessions("ssid", store)
//}
