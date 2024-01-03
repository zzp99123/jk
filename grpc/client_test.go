package grpc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	//初始化一个连接池（准确来说，是池上池）。 d就是一个连接池的池子 就是d里面放了很多个连接池 每一个连接池由ip+端口组成
	//用连接池来初始化一个客户端。
	//利用客户端来发起调用。
	d, err := grpc.Dial(":8090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	c := NewUserServiceClient(d)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	res, err := c.GetById(ctx, &GetByIdRequest{
		Id: 456,
	})
	assert.NoError(t, err)
	t.Log(res.User)
}
