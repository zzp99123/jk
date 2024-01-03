package cache

//
//import (
//	"context"
//	_ "embed"
//	"fmt"
//	"github.com/redis/go-redis/v9"
//	"goFoundation/webook/internal/domain"
//	"strconv"
//	"time"
//)
//
////go:embed lua/interactive.lua
//var luaIncrCnt string
//
//const (
//	fieldReadCnt    = "read_cnt"
//	fieldCollectCnt = "collect_cnt"
//	fieldLikeCnt    = "like_cnt"
//)
//
//type InteractiveCache interface {
//	//收藏 点赞 阅读计数的缓存方法
//	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
//	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
//	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
//	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
//	//查询缓存中数据
//	Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error)
//	//回显数据
//	Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error
//}
//type interactiveCache struct {
//	cmd redis.Cmdable
//}
//
//func NewInteractiveCache(cmd redis.Cmdable) InteractiveCache {
//	return &interactiveCache{
//		cmd: cmd,
//	}
//}
//
//// 阅读计数缓存
//func (i *interactiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
//	return i.cmd.Eval(ctx, luaIncrCnt, []string{i.key(biz, bizId)}, fieldReadCnt, 1).Err()
//}
//
//// 点赞计算缓存
//func (i *interactiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
//	return i.cmd.Eval(ctx, luaIncrCnt, []string{i.key(biz, bizId)}, fieldLikeCnt, 1).Err()
//}
//
//// 取消点赞计数缓存
//func (i *interactiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
//	return i.cmd.Eval(ctx, luaIncrCnt, []string{i.key(biz, bizId)}, fieldLikeCnt, -1).Err()
//}
//
//// 收藏计数缓存
//func (i *interactiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
//	return i.cmd.Eval(ctx, luaIncrCnt, []string{i.key(biz, bizId)}, fieldCollectCnt, 1).Err()
//}
//
//func (i *interactiveCache) Get(ctx context.Context, biz string, bizId int64) (domain.Interactive, error) {
//	//直接使用HGet 即使没有返回值也不会返回err
//	res, err := i.cmd.HGetAll(ctx, i.key(biz, bizId)).Result()
//	if err != nil {
//		return domain.Interactive{}, err
//	}
//	if len(res) == 0 {
//		return domain.Interactive{}, ErrKeyNotExist
//	}
//	readCnt, _ := strconv.ParseInt(res[fieldReadCnt], 10, 64)
//	collectCnt, _ := strconv.ParseInt(res[fieldCollectCnt], 10, 64)
//	likeCnt, _ := strconv.ParseInt(res[fieldLikeCnt], 10, 64)
//	return domain.Interactive{
//		ReadCnt:    readCnt,
//		LikeCnt:    likeCnt,
//		CollectCnt: collectCnt,
//	}, err
//}
//
//func (i *interactiveCache) Set(ctx context.Context, biz string, bizId int64, intr domain.Interactive) error {
//	////TODO implement me
//	//panic("implement me")
//	err := i.cmd.HSet(ctx, i.key(biz, bizId), fieldReadCnt, intr.ReadCnt,
//		fieldCollectCnt, intr.CollectCnt, fieldLikeCnt, intr.LikeCnt).Err()
//	if err != nil {
//		return err
//	}
//	//Redis Expire 命令用于设置 key 的过期时间，key 过期后将不再可用。单位以秒计
//	return i.cmd.Expire(ctx, i.key(biz, bizId), time.Minute*15).Err()
//}
//func (i *interactiveCache) key(biz string, bizId int64) string {
//	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
//}
