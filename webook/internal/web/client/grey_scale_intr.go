// 用适配器和装饰器模式
// 适配InteraciveService，构成一个本质上是本地调用的 gRPC 客户端。
// 将真实的 gRPC 和本地调用的 gRPC 客户端装饰在一起，并且接入流量控制机制
package client

import (
	"context"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	intrv1 "goFoundation/webook/api/proto/gen/intr/v1"
	"google.golang.org/grpc"
	"math/rand"
)

type GreyScaleIntrClinet struct {
	remote    intrv1.InteractiveServiceClient //远程的
	local     intrv1.InteractiveServiceClient //本地的
	threshold *atomicx.Value[int32]
}

func NewGreyScaleIntrClinet(remote intrv1.InteractiveServiceClient, local intrv1.InteractiveServiceClient) *GreyScaleIntrClinet {
	return &GreyScaleIntrClinet{
		remote:    remote,
		local:     local,
		threshold: atomicx.NewValue[int32](),
	}
}
func (g *GreyScaleIntrClinet) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	return g.c().IncrReadCnt(ctx, in, opts...)
}

func (g *GreyScaleIntrClinet) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	return g.c().Like(ctx, in, opts...)
}

func (g *GreyScaleIntrClinet) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	return g.c().CancelLike(ctx, in, opts...)
}

func (g *GreyScaleIntrClinet) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	return g.c().Collect(ctx, in, opts...)
}

func (g *GreyScaleIntrClinet) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	return g.c().Get(ctx, in, opts...)
}

func (g *GreyScaleIntrClinet) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	return g.c().GetByIds(ctx, in, opts...)
}

// 推荐用这种方法 性价比高
func (g *GreyScaleIntrClinet) UpdateThreshold(newThreshold int32) {
	g.threshold.Store(newThreshold) //Store 以原子方式将 val 存储到 x 中。
}

func (g *GreyScaleIntrClinet) OnChange(ch <-chan int32) {
	go func() {
		for res := range ch {
			g.threshold.Store(res)
		}
	}()
}
func (g *GreyScaleIntrClinet) OnChangeV1() chan<- int32 {
	c := make(chan int32, 100)
	go func() {
		for res := range c {
			g.threshold.Store(res)
		}
	}()
	return c
}

// 流量调度算法很简单：
// • 如果随机数小于我们的阈值，就使用远程调用。
// • 否则，使用本地调用。
func (g *GreyScaleIntrClinet) c() intrv1.InteractiveServiceClient {
	res := g.threshold.Load() //Load 以原子方式加载并返回存储在 x 中的值。
	//产生一个随机数。
	num := rand.Int31n(100)
	if num < res {
		return g.remote //远程的
	}
	return g.local //本地
}
