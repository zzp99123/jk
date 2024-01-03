package ioc

import (
	"github.com/spf13/viper"
	"goFoundation/webook/interactive/grpc"
	"goFoundation/webook/pkg/grpcx"
	grpc1 "google.golang.org/grpc"
)

// 如果grpc有多个结构体要初始化 不可能每个都在main函数中写 所以给他单独的拿出来
func InitGrpc(g *grpc.ServerGrpc) *grpcx.Server {
	type config struct {
		Add string
	}
	var c config
	err := viper.UnmarshalKey("grpc.server", &c)
	if err != nil {
		panic(err)
	}
	s := grpc1.NewServer()
	g.Register(s)
	return &grpcx.Server{
		s,
		c.Add,
	}
}
