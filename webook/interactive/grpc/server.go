package grpc

import (
	"context"
	intrv1 "goFoundation/webook/api/proto/gen/intr/v1"
	"goFoundation/webook/interactive/domain"
	"goFoundation/webook/interactive/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ServerGrpc struct {
	intrv1.UnimplementedInteractiveServiceServer
	svc service.InteractiveService
}

func NewServerGrpc(svc service.InteractiveService) *ServerGrpc {
	return &ServerGrpc{
		svc: svc,
	}
}
func (s *ServerGrpc) Register(ser *grpc.Server) {
	intrv1.RegisterInteractiveServiceServer(ser, s)
}
func (s *ServerGrpc) IncrReadCnt(ctx context.Context, request *intrv1.IncrReadCntRequest) (*intrv1.IncrReadCntResponse, error) {
	err := s.svc.IncrReadCnt(ctx, request.GetBiz(), request.GetBizId())
	return &intrv1.IncrReadCntResponse{}, err
}

func (s *ServerGrpc) Like(ctx context.Context, request *intrv1.LikeRequest) (*intrv1.LikeResponse, error) {
	//做校验
	if request.Uid <= 0 {
		return nil, status.Error(codes.InvalidArgument, "uid 错误")
	}
	err := s.svc.Like(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.LikeResponse{}, err
}

func (s *ServerGrpc) CancelLike(ctx context.Context, request *intrv1.CancelLikeRequest) (*intrv1.CancelLikeResponse, error) {
	err := s.svc.CancelLike(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	return &intrv1.CancelLikeResponse{}, err
}

func (s *ServerGrpc) Collect(ctx context.Context, request *intrv1.CollectRequest) (*intrv1.CollectResponse, error) {
	err := s.svc.Collect(ctx, request.GetBiz(), request.GetBizId(), request.GetCid(), request.GetUid())
	return &intrv1.CollectResponse{}, err
}

func (s *ServerGrpc) Get(ctx context.Context, request *intrv1.GetRequest) (*intrv1.GetResponse, error) {
	intr, err := s.svc.Get(ctx, request.GetBiz(), request.GetBizId(), request.GetUid())
	if err != nil {
		return nil, err
	}
	return &intrv1.GetResponse{
		Intr: s.toGrpc(intr),
	}, nil
}

func (s *ServerGrpc) GetByIds(ctx context.Context, request *intrv1.GetByIdsRequest) (*intrv1.GetByIdsResponse, error) {
	intr, err := s.svc.GetByIds(ctx, request.GetBiz(), request.GetBizIds())
	if err != nil {
		return nil, err
	}
	m := make(map[int64]*intrv1.Interactive, len(intr))
	for i, v := range intr {
		m[i] = s.toGrpc(v)
	}
	return &intrv1.GetByIdsResponse{
		Intr: m,
	}, nil
}
func (s *ServerGrpc) toGrpc(intr domain.Interactive) *intrv1.Interactive {
	return &intrv1.Interactive{
		Biz:        intr.Biz,
		BizId:      intr.BizId,
		CollectCnt: intr.CollectCnt,
		Collected:  intr.Collected,
		LikeCnt:    intr.LikeCnt,
		Liked:      intr.Liked,
		ReadCnt:    intr.ReadCnt,
	}
}

//容错方案 缓存方案 项目中难点 亮点
