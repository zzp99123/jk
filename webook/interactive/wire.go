//go:build wireinject

package main

import (
	"github.com/google/wire"
	"goFoundation/webook/interactive/events"
	"goFoundation/webook/interactive/grpc"
	"goFoundation/webook/interactive/ioc"
	"goFoundation/webook/interactive/repository"
	"goFoundation/webook/interactive/repository/cache"
	"goFoundation/webook/interactive/repository/dao"
	"goFoundation/webook/interactive/service"
)

var InitWire = wire.NewSet(ioc.InitRedis, ioc.InitDB, ioc.InitLogger, ioc.InitKafka)
var InteractiveWire = wire.NewSet(
	dao.NewInteractiveDAO,
	cache.NewInteractiveCache,
	repository.NewInteractiveRepository,
	service.NewInteractiveService,
)

func InitAPP() *App {
	wire.Build(
		InitWire,
		InteractiveWire,
		grpc.NewServerGrpc,
		ioc.InitGrpc,
		ioc.NewConsumers,
		events.NewConsumerEvents,
		wire.Struct(new(App), "*"),
	)
	return new(App)
}
