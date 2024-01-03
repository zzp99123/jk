package service

//
//import (
//	"context"
//	"goFoundation/webook/internal/domain"
//	"goFoundation/webook/internal/repository"
//	"goFoundation/webook/pkg/logger"
//	"golang.org/x/sync/errgroup"
//)
//
////go:generate mockgen -source=./interactive.go -package=svcmocks -destination=mocks/interactive.mock.go InteractiveService
//type InteractiveService interface {
//	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
//	Like(ctx context.Context, biz string, bizId int64, id int64) error
//	CancelLike(ctx context.Context, biz string, bizId int64, id int64) error
//	InsertCollectionBiz(ctx context.Context, biz string, cid, bizId, uid int64) error
//	GetInteractiveInfo(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error)
//	GetByIds(ctx context.Context, biz string, bizId []int64) (map[int64]domain.Interactive, error)
//}
//type interactiveService struct {
//	repo repository.InteractiveRepository
//	l    logger.LoggerV1
//}
//
//func NewInteractiveService(repo repository.InteractiveRepository, l logger.LoggerV1) InteractiveService {
//	return &interactiveService{
//		repo: repo,
//		l:    l,
//	}
//}
//
//// 热搜中获取id
//func (i *interactiveService) GetByIds(ctx context.Context, biz string, bizId []int64) (map[int64]domain.Interactive, error) {
//	res, err := i.repo.GetByIds(ctx, biz, bizId)
//	if err != nil {
//		return nil, err
//	}
//	//定义一个map类型
//	m := make(map[int64]domain.Interactive, len(res))
//	for _, v := range res {
//		m[v.BizId] = v
//	}
//	return m, nil
//}
//
//// 浏览数量
//func (i *interactiveService) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
//	return i.repo.IncrReadCnt(ctx, biz, bizId)
//}
//
//// 点赞
//func (i *interactiveService) Like(ctx context.Context, biz string, bizId int64, id int64) error {
//	return i.repo.Like(ctx, biz, bizId, id)
//}
//
//// 取消点赞
//func (i *interactiveService) CancelLike(ctx context.Context, biz string, bizId int64, id int64) error {
//	return i.repo.CancelLike(ctx, biz, bizId, id)
//}
//
//// 收藏
//func (i *interactiveService) InsertCollectionBiz(ctx context.Context, biz string, cid, bizId, uid int64) error {
//	return i.repo.InsertCollectionBiz(ctx, biz, cid, uid, bizId)
//}
//
//// 搜素获取收藏的信息
//func (i *interactiveService) GetInteractiveInfo(ctx context.Context, biz string, bizId, uid int64) (domain.Interactive, error) {
//	res, err := i.repo.GetInteractiveInfo(ctx, biz, bizId)
//	if err != nil {
//		return domain.Interactive{}, err
//	}
//	var eg errgroup.Group
//
//	eg.Go(func() error {
//		res.Liked, err = i.repo.GetLikeInfo(ctx, biz, bizId, uid)
//		return err
//	})
//	eg.Go(func() error {
//		res.Collected, err = i.repo.GetCollectionInfo(ctx, biz, bizId, uid)
//		return err
//	})
//	// 说明是登录过的，补充用户是否点赞或者
//	// 新的打印日志的形态 zap 本身就有这种用法
//	//上面的return err 返回到这里 然后做日志
//	err = eg.Wait()
//	if err != nil {
//		i.l.Error("查询用户是否点赞的信息失败",
//			logger.String("biz", biz),
//			logger.Int64("bizId", bizId),
//			logger.Int64("uid", uid),
//			logger.Error(err))
//	}
//	return res, err
//}
