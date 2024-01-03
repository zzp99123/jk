package article

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

type daoArticle struct {
	db *gorm.DB
}

func NewDaoArticle(db *gorm.DB) DaoArticle {
	return &daoArticle{
		db: db,
	}
}

// 热搜
func (dao *daoArticle) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]Article, error) {
	var a []Article
	err := dao.db.WithContext(ctx).Where("utime < ?", start.UnixMilli()).
		Order("utime DESC").Offset(offset).Limit(limit).Find(&a).Error
	return a, err
}

// 创作者查id
func (dao *daoArticle) Detail(ctx context.Context, id int64) (Article, error) {
	var art Article
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&art).Error
	return art, err
}

// 读者查id
func (dao *daoArticle) PubDetail(ctx context.Context, id int64) (PublishedArticle, error) {
	var pub PublishedArticle
	err := dao.db.WithContext(ctx).Where("id=?", id).First(&pub).Error
	return pub, err
}

// 搜素
func (dao *daoArticle) List(ctx context.Context, id int64, Offset, Limit int) ([]Article, error) {
	var a []Article
	err := dao.db.WithContext(ctx).Model(&Article{}).
		Where("author_id=?", id).
		Offset(Offset).
		Limit(Limit).
		// 升序排序。 utime ASC
		// 混合排序
		// ctime ASC, utime desc
		Order("utime DESC").
		//第二种写法
		//Order(clause.OrderBy{Columns: []clause.OrderByColumn{
		//	{Column: clause.Column{Name: "utime"}, Desc: true},
		//	{Column: clause.Column{Name: "ctime"}, Desc: false},
		//}}).
		First(&a).Error
	return a, err
}

// 事务
func (dao *daoArticle) Transaction(ctx context.Context,
	bizFunc func(txDAO DaoArticle) error) error {
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewDaoArticle(tx)
		return bizFunc(txDAO)
	})
}

// 仅自己可见
func (dao *daoArticle) Withdraw(ctx context.Context, id int64, author int64, status uint8) error {
	//仅自己可见需要改变状态 所以用到事务 要把制作库和线上库都改变
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).
			Where("id=? AND author_id=?", id, author).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("用户在尝试操作非本人数据")
		}
		res = tx.Model(&PublishedArticle{}).
			Where("id=? AND author_id=?", id, author).
			Updates(map[string]any{
				"status": status,
				"utime":  now,
			})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return errors.New("用户在尝试操作非本人数据")
		}
		return nil
	})
}

// 发表
func (dao *daoArticle) Sync(ctx context.Context, art Article) (int64, error) {
	tx := dao.db.WithContext(ctx).Begin()
	now := time.Now().UnixMilli()
	defer tx.Rollback()
	txDAO := NewDaoArticle(tx)
	var (
		id  = art.Id
		err error
	)
	if id == 0 {
		id, err = txDAO.Create(ctx, art)
	} else {
		err = txDAO.Update(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	art.Id = id
	publishArt := PublishedArticle(art)
	publishArt.Utime = now
	publishArt.Ctime = now
	err = tx.Clauses(clause.OnConflict{
		// ID 冲突的时候。实际上，在 MYSQL 里面你写不写都可以
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"title":   art.Title,
			"content": art.Content,
			"status":  art.Status,
			"utime":   now,
		}),
	}).Create(&publishArt).Error
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return id, tx.Error
}

// 创建
func (a *daoArticle) Create(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	err := a.db.WithContext(ctx).Create(&art).Error
	return art.Id, err
}

// 更改
func (a *daoArticle) Update(ctx context.Context, art Article) error {
	now := time.Now().UnixMilli()
	art.Utime = now
	res := a.db.WithContext(ctx).Model(&art).Where("id=? AND author_id=?", art.Id, art.AuthorId).
		Updates(map[string]any{
			"title":   art.Title,
			"content": art.Content,
			"utime":   art.Utime,
			"status":  art.Status,
		})
	if res.Error != nil {
		return res.Error
	}
	//获取更新的条数
	if res.RowsAffected == 0 {
		return fmt.Errorf("更新失败，可能是创作者非法 id %d, author_id %d",
			art.Id, art.AuthorId)
	}
	return res.Error
}
