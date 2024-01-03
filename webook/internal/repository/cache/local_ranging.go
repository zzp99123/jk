// 本地缓存
package cache

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/syncx/atomicx"
	"goFoundation/webook/internal/domain"
	"time"
)

type localRanging interface {
	Set(ctx context.Context, a []domain.Article) error
	Get(ctx context.Context) ([]domain.Article, error)
	ForceGet(ctx context.Context) ([]domain.Article, error)
}
type LocalRanging struct {
	// 我用我的泛型封装
	// 你可以考虑直接使用 uber 的，或者 SDK 自带的
	art        *atomicx.Value[[]domain.Article]
	ddl        *atomicx.Value[time.Time]
	expiration time.Duration
}

func NewLocalRanging() *LocalRanging {
	return &LocalRanging{
		art: atomicx.NewValue[[]domain.Article](),
		ddl: atomicx.NewValueOf(time.Now()),
		// 永不过期，或者非常长，或者对齐到 redis 的过期时间，都行
		expiration: time.Minute * 10,
	}
}
func (l *LocalRanging) Set(ctx context.Context, a []domain.Article) error {
	l.art.Store(a)
	ddl := time.Now().Add(l.expiration)
	l.ddl.Store(ddl)
	return nil
}
func (l *LocalRanging) Get(ctx context.Context) ([]domain.Article, error) {
	ddl := l.ddl.Load()
	art := l.art.Load()
	if len(art) == 0 || ddl.Before(time.Now()) {
		return nil, errors.New("本地缓存未命中")
	}
	return art, nil
}
func (l *LocalRanging) ForceGet(ctx context.Context) ([]domain.Article, error) {
	arts := l.art.Load()
	return arts, nil
}

type item struct {
	arts []domain.Article
	ddl  time.Time
}
