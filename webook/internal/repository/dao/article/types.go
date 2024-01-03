package article

import (
	"context"
	"time"
)

type DaoArticle interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
	Sync(ctx context.Context, art Article) (int64, error)
	Withdraw(ctx context.Context, id int64, author int64, status uint8) error
	List(ctx context.Context, id int64, Offset, Limit int) ([]Article, error)
	Detail(ctx context.Context, id int64) (Article, error)
	PubDetail(ctx context.Context, id int64) (PublishedArticle, error)
	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]Article, error)
}
