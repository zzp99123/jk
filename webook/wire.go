//go:build wireinject

package main

import (
	"github.com/google/wire"
	events1 "goFoundation/webook/interactive/events"
	repository2 "goFoundation/webook/interactive/repository"
	cache2 "goFoundation/webook/interactive/repository/cache"
	dao2 "goFoundation/webook/interactive/repository/dao"
	service2 "goFoundation/webook/interactive/service"
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

// 用户邮箱密码注册登录
var UserWire = wire.NewSet(
	dao.NewUserDao,
	cache.NewUsersCache,
	repository.NewUserRepository,
	service.NewUserService,
	web.NewUserHandle,
)

// 验证码登录
var CodeWire = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	service.NewServiceCode,
)

// 文章
var ArticleWire = wire.NewSet(
	article2.NewDaoArticle,
	service.NewArticleService,
	article.NewArticleRepository,
	cache.NewArticleCache,
	web.NewArticleHandlers,
)

// 阅读，收藏，点赞
var InteractiveWire = wire.NewSet(
	dao2.NewInteractiveDAO,
	cache2.NewInteractiveCache,
	repository2.NewInteractiveRepository,
	service2.NewInteractiveService,
)

// 文章热搜
var RangWire = wire.NewSet(
	cache.NewRangingCache,
	cache.NewLocalRanging,
	repository.NewRangingRepository,
	service.NewRangingService,
)

func InitWebServer() *App {
	wire.Build(
		//
		ioc.InitDB, ioc.InitRedis, ioc.InitWechatService,
		ioc.InitSMSService, ioc.InitWebServer, ioc.InitMiddlewares,
		ioc.NewWechatHandlerConfig,
		ioc.InitKafka, ioc.NewConsumers, ioc.NewSyncProducer, ioc.InitLogger,
		ioc.InitJobs, ioc.InitRankingJob, ioc.InitRLockClient, ioc.InitIntrGRPCClient,

		//jwt缓存
		ijwt.NewRedisJWTHandler,
		//微信登录
		web.NewOAuth2WechatHandler,
		//kafka
		//events.NewConsumerEvents,
		events1.NewInteractiveReadEventBatchConsumer,
		events.NewProducerEvents,
		UserWire,
		CodeWire,
		ArticleWire,
		InteractiveWire,
		RangWire,
		// 组装我这个结构体的所有字段
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
