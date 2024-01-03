//go:build wireinject

package startup

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	events "goFoundation/webook/internal/events/article"
	"goFoundation/webook/internal/repository"
	"goFoundation/webook/internal/repository/article"
	"goFoundation/webook/internal/repository/cache"
	"goFoundation/webook/internal/repository/dao"
	article2 "goFoundation/webook/internal/repository/dao/article"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/internal/web"
	ijwt "goFoundation/webook/internal/web/jwt"
	"goFoundation/webook/ioc"
)

var thirdProvider = wire.NewSet(InitRedis, InitTestDB, InitLog, ioc.NewSyncProducer, ioc.InitKafka)
var userSvcProvider = wire.NewSet(
	dao.NewUserDao,
	cache.NewUsersCache,
	repository.NewUserRepository,
	service.NewUserService)
var CodeProvider = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	service.NewServiceCode,
)
var articleSvcProvider = wire.NewSet(
	article2.NewDaoArticle,
	service.NewArticleService,
	article.NewArticleRepository,
	cache.NewArticleCache,
)

func InitWebServer() *gin.Engine {
	wire.Build(
		thirdProvider,
		ioc.InitSMSService,
		// 指定啥也不干的 wechat service
		InitPhantomWechatService,
		events.NewProducerEvents,
		// handler 部分
		web.NewUserHandle,
		web.NewOAuth2WechatHandler,
		web.NewArticleHandlers,
		InitWechatHandlerConfig,
		ijwt.NewRedisJWTHandler,
		// gin 的中间件
		ioc.InitMiddlewares,
		// Web 服务器
		ioc.InitWebServer,
		userSvcProvider,
		CodeProvider,
		articleSvcProvider,
	)
	// 随便返回一个
	return gin.Default()
}

func InitArticleHandler(dao article2.DaoArticle) *web.ArticleHandler {
	wire.Build(thirdProvider,
		userSvcProvider,
		events.NewProducerEvents,
		service.NewArticleService,
		article.NewArticleRepository,
		cache.NewArticleCache,
		web.NewArticleHandlers,
	)
	return &web.ArticleHandler{}
}

func InitUserSvc() service.UserService {
	wire.Build(thirdProvider, userSvcProvider)
	return service.NewUserService(nil, nil)
}

func InitJwtHdl() ijwt.Handler {
	wire.Build(thirdProvider, ijwt.NewRedisJWTHandler)
	return ijwt.NewRedisJWTHandler(nil)
}
