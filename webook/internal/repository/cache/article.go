package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goFoundation/webook/internal/domain"
	"time"
)

type ArticleCache interface {
	// GetFirstPage 只缓存第第一页的数据
	// 并且不缓存整个 Content
	GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error)
	SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error
	DelFirstPage(ctx context.Context, author int64) error
	//创作者的缓存
	Set(ctx context.Context, art domain.Article) error
	Get(ctx context.Context, id int64) (domain.Article, error)
	//读者的缓存
	SetPub(ctx context.Context, art domain.Article) error
	GetPub(ctx context.Context, id int64) (domain.Article, error)
}
type articleCache struct {
	cmd redis.Cmdable
}

func NewArticleCache(cmd redis.Cmdable) ArticleCache {
	return &articleCache{
		cmd: cmd,
	}
}

func (a *articleCache) GetFirstPage(ctx context.Context, author int64) ([]domain.Article, error) {
	res, err := a.cmd.Get(ctx, a.firstPageKey(author)).Bytes()
	if err != nil {
		return nil, err
	}
	var art []domain.Article
	err = json.Unmarshal(res, &art)
	return art, err

}

func (a *articleCache) SetFirstPage(ctx context.Context, author int64, arts []domain.Article) error {
	for i := range arts {
		arts[i].Content = arts[i].Abstract()
	}
	res, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	err = a.cmd.Set(ctx, a.firstPageKey(author), res, time.Minute*10).Err()
	return err
}

func (a *articleCache) DelFirstPage(ctx context.Context, author int64) error {
	return a.cmd.Del(ctx, a.firstPageKey(author)).Err()
}

func (a *articleCache) Set(ctx context.Context, art domain.Article) error {
	res, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.cmd.Set(ctx, a.authorArtKey(art.Id), res, time.Minute).Err()
}

func (a *articleCache) Get(ctx context.Context, id int64) (domain.Article, error) {
	res, err := a.cmd.Get(ctx, a.authorArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(res, &art)
	return art, err
}
func (a *articleCache) SetPub(ctx context.Context, art domain.Article) error {
	res, err := json.Marshal(art)
	if err != nil {
		return err
	}
	return a.cmd.Set(ctx, a.readerArtKey(art.Id), res, time.Minute).Err()
}

func (a *articleCache) GetPub(ctx context.Context, id int64) (domain.Article, error) {
	res, err := a.cmd.Get(ctx, a.readerArtKey(id)).Bytes()
	if err != nil {
		return domain.Article{}, err
	}
	var art domain.Article
	err = json.Unmarshal(res, &art)
	return art, err
}
func (a *articleCache) firstPageKey(id int64) string {
	return fmt.Sprintf("article:first_page:%d", id)
}

// 创作者
func (a *articleCache) authorArtKey(id int64) string {
	return fmt.Sprintf("article:author:%d", id)
}

// 读者
func (a *articleCache) readerArtKey(id int64) string {
	return fmt.Sprintf("article:reader:%d", id)
}
