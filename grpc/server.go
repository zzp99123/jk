package grpc

import "context"

type server struct {
	UnimplementedUserServiceServer
}

var _UserServiceServer = &server{}

func (s *server) GetById(ctx context.Context, g *GetByIdRequest) (*GetByIdResponse, error) {
	return &GetByIdResponse{
		User: &User{
			Id:   2,
			Name: "abcd",
		},
	}, nil
}
