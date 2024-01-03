package grpcx

import (
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	*grpc.Server
	Addr string
}

func (s *Server) Serve() error {
	// 监听 8090 端口，你可以随便写
	l, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	return s.Server.Serve(l)
}
