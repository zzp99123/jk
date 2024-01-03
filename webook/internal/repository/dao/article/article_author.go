package article

import (
	"context"
	"gorm.io/gorm"
)

type DaoArticleAuthor interface {
	Create(ctx context.Context, art Article) (int64, error)
	Update(ctx context.Context, art Article) error
}

func NewDaoArticleAuthor(db *gorm.DB) DaoArticleAuthor {
	panic("implement me")
}

//func (a *daoArticleAuthor) Create(ctx context.Context, art Article) (int64, error) {
//	now := time.Now().UnixMilli()
//	art.Ctime = now
//	art.Utime = now
//	err := a.db.WithContext(ctx).Create(&art).Error
//	return art.Id, err
//}
//
//func (a *daoArticleAuthor) Update(ctx context.Context, art Article) error {
//	now := time.Now().UnixMilli()
//	art.Utime = now
//	res := a.db.WithContext(ctx).Model(&art).Where("id=? AND author_id=?", art.Id, art.AuthorId).
//		Updates(map[string]any{
//			"title":   art.Title,
//			"content": art.Content,
//			"utime":   art.Utime,
//		})
//	if res.Error != nil {
//		return res.Error
//	}
//	//获取更新的条数
//	if res.RowsAffected == 0 {
//		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d",
//			art.Id, art.AuthorId)
//	}
//	return res.Error
//}
