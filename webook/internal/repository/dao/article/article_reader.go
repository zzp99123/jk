package article

import (
	"context"
	"gorm.io/gorm"
)

type DaoArticleReader interface {
	Save(ctx context.Context, art Article) error
	SaveV2(ctx context.Context, art PublishedArticle) error
}

func NewDaoArticleReader(db *gorm.DB) DaoArticleReader {
	panic("implement me")
}

//func (d *PublishedArticle) Save(ctx context.Context, art Article) error {
//	//return nil
//}
