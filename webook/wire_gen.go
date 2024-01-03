// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"goFoundation/webook/interactive/events"
	repository2 "goFoundation/webook/interactive/repository"
	cache2 "goFoundation/webook/interactive/repository/cache"
	dao2 "goFoundation/webook/interactive/repository/dao"
	service2 "goFoundation/webook/interactive/service"
	article3 "goFoundation/webook/internal/events/article"
	"goFoundation/webook/internal/repository"
	article2 "goFoundation/webook/internal/repository/article"
	"goFoundation/webook/internal/repository/cache"
	"goFoundation/webook/internal/repository/dao"
	"goFoundation/webook/internal/repository/dao/article"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/internal/web"
	"goFoundation/webook/internal/web/jwt"
	"goFoundation/webook/ioc"
)

// Injectors from wire.go:

func InitWebServer() *App {
	cmdable := ioc.InitRedis()
	handler := jwt.NewRedisJWTHandler(cmdable)
	loggerV1 := ioc.InitLogger()
	v := ioc.InitMiddlewares(cmdable, handler, loggerV1)
	db := ioc.InitDB(loggerV1)
	userDao := dao.NewUserDao(db)
	usersCache := cache.NewUsersCache(cmdable)
	userRepository := repository.NewUserRepository(userDao, usersCache)
	userService := service.NewUserService(userRepository, loggerV1)
	codeCache := cache.NewCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService()
	codeService := service.NewServiceCode(codeRepository, smsService)
	userHandler := web.NewUserHandle(userService, codeService, handler)
	wechatService := ioc.InitWechatService(loggerV1)
	wechatHandlerConfig := ioc.NewWechatHandlerConfig()
	oAuth2WechatHandler := web.NewOAuth2WechatHandler(wechatService, userService, wechatHandlerConfig, handler)
	daoArticle := article.NewDaoArticle(db)
	articleCache := cache.NewArticleCache(cmdable)
	articleRepository := article2.NewArticleRepository(daoArticle, userRepository, articleCache, loggerV1)
	client := ioc.InitKafka()
	syncProducer := ioc.NewSyncProducer(client)
	producerEvents := article3.NewProducerEvents(syncProducer)
	articleService := service.NewArticleService(articleRepository, loggerV1, producerEvents)
	interactiveDAO := dao2.NewInteractiveDAO(db)
	interactiveCache := cache2.NewInteractiveCache(cmdable)
	interactiveRepository := repository2.NewInteractiveRepository(interactiveDAO, interactiveCache, loggerV1)
	interactiveService := service2.NewInteractiveService(interactiveRepository, loggerV1)
	interactiveServiceClient := ioc.InitIntrGRPCClient(interactiveService)
	articleHandler := web.NewArticleHandlers(articleService, interactiveServiceClient, loggerV1)
	engine := ioc.InitWebServer(v, userHandler, oAuth2WechatHandler, articleHandler)
	interactiveReadEventBatchConsumer := events.NewInteractiveReadEventBatchConsumer(client, loggerV1, interactiveRepository)
	v2 := ioc.NewConsumers(interactiveReadEventBatchConsumer)
	rangingService := service.NewRangingService(articleService, interactiveServiceClient)
	rlockClient := ioc.InitRLockClient(cmdable)
	rankingJob := ioc.InitRankingJob(rangingService, rlockClient, loggerV1)
	cron := ioc.InitJobs(loggerV1, rankingJob)
	app := &App{
		web:       engine,
		consumers: v2,
		cron:      cron,
	}
	return app
}

// wire.go:

// 用户邮箱密码注册登录
var UserWire = wire.NewSet(dao.NewUserDao, cache.NewUsersCache, repository.NewUserRepository, service.NewUserService, web.NewUserHandle)

// 验证码登录
var CodeWire = wire.NewSet(cache.NewCodeCache, repository.NewCodeRepository, service.NewServiceCode)

// 文章
var ArticleWire = wire.NewSet(article.NewDaoArticle, service.NewArticleService, article2.NewArticleRepository, cache.NewArticleCache, web.NewArticleHandlers)

// 阅读，收藏，点赞
var InteractiveWire = wire.NewSet(dao2.NewInteractiveDAO, cache2.NewInteractiveCache, repository2.NewInteractiveRepository, service2.NewInteractiveService)

// 文章热搜
var RangWire = wire.NewSet(cache.NewRangingCache, cache.NewLocalRanging, repository.NewRangingRepository, service.NewRangingService)
