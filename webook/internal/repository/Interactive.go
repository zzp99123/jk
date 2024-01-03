package repository

//
//import (
//	"context"
//	"github.com/ecodeclub/ekit/slice"
//	"goFoundation/webook/internal/domain"
//	"goFoundation/webook/internal/repository/cache"
//	"goFoundation/webook/internal/repository/dao"
//	"goFoundation/webook/pkg/logger"
//)
//
//type InteractiveRepository interface {
//	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
//	BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error
//	Like(ctx context.Context, biz string, bizId int64, id int64) error
//	CancelLike(ctx context.Context, biz string, bizId int64, id int64) error
//	InsertCollectionBiz(ctx context.Context, biz string, cid, bizId, uid int64) error
//	GetInteractiveInfo(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
//	GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (bool, error)
//	GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (bool, error)
//	AddRecord(ctx context.Context, aid int64, uid int64) error
//	GetByIds(ctx context.Context, biz string, bizId []int64) ([]domain.Interactive, error)
//}
//type interactiveRepository struct {
//	dao        dao.InteractiveDao
//	interCache cache.InteractiveCache
//	l          logger.LoggerV1
//}
//
//func NewInteractiveRepository(dao dao.InteractiveDao, interCache cache.InteractiveCache, l logger.LoggerV1) InteractiveRepository {
//	return &interactiveRepository{
//		dao:        dao,
//		interCache: interCache,
//		l:          l,
//	}
//}
//
//// 热搜中获取id
//func (i *interactiveRepository) GetByIds(ctx context.Context, biz string, bizId []int64) ([]domain.Interactive, error) {
//	res, err := i.dao.GetByIds(ctx, biz, bizId)
//	if err != nil {
//		return nil, err
//	}
//	return slice.Map(res, func(idx int, src dao.Interactive) domain.Interactive {
//		return i.toDomain(src)
//	}), nil
//
//}
//
//// 阅读记录功能 就是简单的增删改查
//func (i *interactiveRepository) AddRecord(ctx context.Context, aid int64, uid int64) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//// 批量浏览数量
//func (i *interactiveRepository) BatchIncrReadCnt(ctx context.Context, biz []string, bizId []int64) error {
//	err := i.dao.BatchIncrReadCnt(ctx, biz, bizId)
//	if err != nil {
//		return err
//	}
//	// 你也要批量的去修改 redis，所以就要去改 lua 脚本
//	// c.cache.IncrReadCntIfPresent()
//	// TODO, 等我写新的 lua 脚本/或者用 pipeline
//	return nil
//}
//
//// 浏览数量
//func (i *interactiveRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
//	//先查数据库 在查缓存 因为缓存里没有直接查key的方法 所以 你得用HIncrBy这个命令 但是这个命令的缺陷就是 如果 read_cnt 不存在，
//	//就先设置为 0，而后自增 1 如果数据库已经存了10了 你如果先查缓存的话 他就变成1了 所以 得得先用数据库
//	err := i.dao.IncrReadCnt(ctx, biz, bizId)
//	if err != nil {
//		return err
//	}
//	return i.interCache.IncrReadCntIfPresent(ctx, biz, bizId)
//}
//
//// 点赞
//func (i *interactiveRepository) Like(ctx context.Context, biz string, bizId int64, id int64) error {
//	err := i.dao.Like(ctx, biz, bizId, id)
//	if err != nil {
//		return err
//	}
//	return i.interCache.IncrLikeCntIfPresent(ctx, biz, bizId)
//}
//
//// 取消点赞
//func (i *interactiveRepository) CancelLike(ctx context.Context, biz string, bizId int64, id int64) error {
//	err := i.dao.CancelLike(ctx, biz, bizId, id)
//	if err != nil {
//		return err
//	}
//	return i.interCache.DecrLikeCntIfPresent(ctx, biz, bizId)
//}
//
//// 收藏
//func (i *interactiveRepository) InsertCollectionBiz(ctx context.Context, biz string, cid, bizId, uid int64) error {
//	//数据库
//	err := i.dao.InsertCollectionBiz(ctx, dao.UserCollectionBiz{
//		BizId: bizId,
//		Biz:   biz,
//		Uid:   uid,
//		Cid:   cid,
//	})
//	if err != nil {
//		return err
//	}
//	//缓存
//	return i.interCache.IncrCollectCntIfPresent(ctx, biz, bizId)
//}
//
//// 获取用户的收藏，点赞，阅读的总数 并且回显
//func (i *interactiveRepository) GetInteractiveInfo(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
//	//查询缓存中的信息
//	data, err := i.interCache.Get(ctx, biz, bizId)
//	if err == nil {
//		// 缓存只缓存了具体的数字，但是没有缓存自身有没有点赞的信息
//		// 因为一个人反复刷，重复刷一篇文章是小概率的事情
//		// 也就是说，你缓存了某个用户是否点赞的数据，命中率会很低
//		return data, err
//	}
//	inter, err := i.dao.GetInteractiveInfo(ctx, biz, bizId)
//	if err == nil {
//		res := i.toDomain(inter)
//		er := i.interCache.Set(ctx, biz, bizId, res)
//		if er != nil {
//			i.l.Error("回写缓存失败",
//				logger.Int64("bizId", bizId),
//				logger.String("biz", biz),
//				logger.Error(er))
//		}
//	}
//	return domain.Interactive{}, err
//}
//
//// 获取收藏用户的信息
//func (i *interactiveRepository) GetCollectionInfo(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
//	_, err := i.dao.GetCollectionInfo(ctx, biz, bizId, uid)
//	switch err {
//	case nil:
//		return true, nil
//	case dao.ErrDataNotFound:
//		return false, nil
//	default:
//		return false, err
//	}
//}
//func (i *interactiveRepository) GetLikeInfo(ctx context.Context, biz string, bizId, uid int64) (bool, error) {
//	_, err := i.dao.GetLikeInfo(ctx, biz, bizId, uid)
//	switch err {
//	case nil:
//		return true, nil
//	case dao.ErrDataNotFound:
//		return false, nil
//	default:
//		return false, err
//	}
//}
//func (i *interactiveRepository) toDomain(art dao.Interactive) domain.Interactive {
//	return domain.Interactive{
//		ReadCnt:    art.ReadCnt,
//		LikeCnt:    art.LikeCnt,
//		CollectCnt: art.CollectCnt,
//	}
//}
