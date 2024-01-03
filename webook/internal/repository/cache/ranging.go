package cache

import (
	"context"
	"encoding/json"
	"github.com/redis/go-redis/v9"
	"goFoundation/webook/internal/domain"
	"time"
)

type rangingCache interface {
	Get(ctx context.Context) ([]domain.Article, error)
	Set(ctx context.Context, a []domain.Article) error
}
type RangingCache struct {
	cmd      redis.Cmdable
	duration time.Duration
	key      string
	loading  string
}

func NewRangingCache(cmd redis.Cmdable) *RangingCache {
	return &RangingCache{
		cmd: cmd,
	}
}
func (r *RangingCache) Get(ctx context.Context) ([]domain.Article, error) {
	res, err := r.cmd.Get(ctx, r.key).Bytes()
	if err != nil {
		return nil, err
	}
	var a []domain.Article
	err = json.Unmarshal(res, &a)
	return a, err

}
func (r *RangingCache) Set(ctx context.Context, a []domain.Article) error {
	// 你可以趁机，把 article 写到缓存里面 id => article
	for i := 0; i < len(a); i++ {
		a[i].Content = ""
	}
	val, err := json.Marshal(a)
	if err != nil {
		return err
	}
	// 这个过期时间要稍微长一点，最好是超过计算热榜的时间（包含重试在内的时间）
	// 你甚至可以直接永不过期
	return r.cmd.Set(ctx, r.key, val, r.duration).Err()
}

//我刚进公司的时候只有redis缓存，随着业务量增大，性能太差
//我就引用了本地缓存，后来我再一次改进，我考虑到分布式环境，本地缓存订阅更新的技术，让所有实例统一更新他的本地缓存
//再然后我又想到既然我的实例可以这样 那么我可以协调bff那边 让他们做他们那边的本地缓存
