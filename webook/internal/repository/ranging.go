package repository

import (
	"context"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository/cache"
)

// 热度只需要存在缓存里 不用存在数据库
type RangingRepository interface {
	TopN(ctx context.Context, a []domain.Article) error
	GetTopN(ctx context.Context) ([]domain.Article, error)
}

type rangingRepository struct {
	Cache      *cache.RangingCache
	localCache *cache.LocalRanging
}

func NewRangingRepository(Cache *cache.RangingCache, localCache *cache.LocalRanging) RangingRepository {
	return &rangingRepository{
		Cache:      Cache,
		localCache: localCache,
	}
}
func (r *rangingRepository) TopN(ctx context.Context, a []domain.Article) error {
	// 这一步必然不会出错
	_ = r.localCache.Set(ctx, a)
	return r.Cache.Set(ctx, a)
}
func (r *rangingRepository) GetTopN(ctx context.Context) ([]domain.Article, error) {
	art, err := r.localCache.Get(ctx)
	if err != nil {
		return nil, err
	}
	art, err = r.Cache.Get(ctx)
	if err == nil {
		_ = r.Cache.Set(ctx, art)
	} else {
		// 这里，我们没有进一步区分是什么原因导致的 Redis 错误
		return r.localCache.ForceGet(ctx)
	}
	return art, nil
}
