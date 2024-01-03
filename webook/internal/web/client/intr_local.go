//我们的最开始的方案里面为了保证迁移顺利，并且保证可以回滚，我们有一个阶段是请求可以是本地调用 InteraciveService 的，
//也可以是真的发起了 gRPC 调用，这个过程我们需要进一步控制住流量。

package client

import (
	"context"
	intrv1 "goFoundation/webook/api/proto/gen/intr/v1"
	"goFoundation/webook/interactive/domain"
	"goFoundation/webook/interactive/service"
	"google.golang.org/grpc"
)

type IntrLocalClient struct {
	svc service.InteractiveService
}

func NewIntrLocalClient(svc service.InteractiveService) *IntrLocalClient {
	return &IntrLocalClient{
		svc: svc,
	}
}

func (i *IntrLocalClient) IncrReadCnt(ctx context.Context, in *intrv1.IncrReadCntRequest, opts ...grpc.CallOption) (*intrv1.IncrReadCntResponse, error) {
	err := i.svc.IncrReadCnt(ctx, in.GetBiz(), in.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

func (i *IntrLocalClient) Like(ctx context.Context, in *intrv1.LikeRequest, opts ...grpc.CallOption) (*intrv1.LikeResponse, error) {
	err := i.svc.Like(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.LikeResponse{}, err
}

func (i *IntrLocalClient) CancelLike(ctx context.Context, in *intrv1.CancelLikeRequest, opts ...grpc.CallOption) (*intrv1.CancelLikeResponse, error) {
	err := i.svc.CancelLike(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	return &intrv1.CancelLikeResponse{}, err
}

func (i *IntrLocalClient) Collect(ctx context.Context, in *intrv1.CollectRequest, opts ...grpc.CallOption) (*intrv1.CollectResponse, error) {
	err := i.svc.Collect(ctx, in.GetBiz(), in.GetBizId(), in.GetCid(), in.GetUid())
	return &intrv1.CollectResponse{}, err
}

func (i *IntrLocalClient) Get(ctx context.Context, in *intrv1.GetRequest, opts ...grpc.CallOption) (*intrv1.GetResponse, error) {
	res, err := i.svc.Get(ctx, in.GetBiz(), in.GetBizId(), in.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: i.ToGrpc(res),
	}, nil
}

func (i *IntrLocalClient) GetByIds(ctx context.Context, in *intrv1.GetByIdsRequest, opts ...grpc.CallOption) (*intrv1.GetByIdsResponse, error) {
	res, err := i.svc.GetByIds(ctx, in.GetBiz(), in.GetBizIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*intrv1.Interactive, len(res))
	for c, v := range res {
		m[c] = i.ToGrpc(v)
	}
	return &intrv1.GetByIdsResponse{
		Intr: m,
	}, nil
}
func (i *IntrLocalClient) ToGrpc(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		LikeCnt:    intr.LikeCnt,
		Liked:      intr.Liked,
		ReadCnt:    intr.ReadCnt}
}
