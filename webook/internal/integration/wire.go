//go:build wireinject

package integration

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"goFoundation/webook/internal/repository"
	"goFoundation/webook/internal/repository/cache"
	"goFoundation/webook/internal/repository/dao"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/internal/web"
	"goFoundation/webook/ioc"
)

func InitWebServer() *gin.Engine {
	wire.Build(
		//
		ioc.InitDb, ioc.InitRedis,
		dao.NewUserService, cache.NewCacheUsers, cache.NewCacheCode,
		repository.NewUserService, repository.NewRepositoryCode,
		service.NewUserService, service.NewServiceCode, ioc.InitSMSService,
		web.NewUserHandle, ioc.InitWebServer, ioc.InitMiddlewares,
	)
	return new(gin.Engine)
}
