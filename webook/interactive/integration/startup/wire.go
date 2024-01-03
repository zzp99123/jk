//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"goFoundation/webook/interactive/grpc"
	"goFoundation/webook/interactive/repository"
	"goFoundation/webook/interactive/repository/cache"
	"goFoundation/webook/interactive/repository/dao"
	"goFoundation/webook/interactive/service"
)

var InitWire = wire.NewSet(InitTestDB, InitRedis, InitLog)
var InteractiveWire = wire.NewSet(
	dao.NewInteractiveDAO,
	cache.NewInteractiveCache,
	repository.NewInteractiveRepository,
	service.NewInteractiveService,
)

func InitInteractiveService() service.InteractiveService {
	wire.Build(InitWire, InteractiveWire)
	return service.NewInteractiveService(nil, nil)
}
func InitInteractiveGRPCServer() *grpc.ServerGrpc {
	wire.Build(InitWire, InteractiveWire, grpc.NewServerGrpc)
	return new(grpc.ServerGrpc)
}
