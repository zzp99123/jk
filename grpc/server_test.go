package grpc

import (
	"google.golang.org/grpc"
	"net"
	"testing"
)

func TestServer(t *testing.T) {
	//先创建一个 gRPC Server
	s := grpc.NewServer()
	defer func() {
		// 优雅退出
		s.GracefulStop()
	}()
	//再创建一个 UserServiceServer 实现的实例，
	us := &server{}
	//调用RegisterUserServiceServer 注册一下。这个方法是Protobuf 生成的
	RegisterUserServiceServer(s, us)
	//创建一个监听网络端口的 Listener。
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	//调用 gRPC Server 上的 Serve 方法。
	err = s.Serve(l)
	t.Log(err)
}
