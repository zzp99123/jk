// 线上存档用于大文本数据存储  线上库 可以存mongo也可以用这个oss
package article

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ecodeclub/ekit"
	"goFoundation/webook/internal/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"strconv"
	"time"
)

var statusPrivate = domain.ArticleStatusPrivate.ToUint8()

type S3DAO struct {
	oss    *s3.S3
	bucket *string
	daoArticle
	// 通过组合 daoArticle 来简化操作
	// 当然在实践中，你是不太会有组合的机会
	// 你操作制作库总是一样的
	// 你就是操作线上库的时候不一样
}

// NewOssDAO 因为组合 daoArticle 是一个内部实现细节
// 所以这里要直接传入 DB
func NewOssDAO(oss *s3.S3, db *gorm.DB) DaoArticle {
	return &S3DAO{
		oss: oss,
		// 你也可以考虑利用依赖注入来传入。
		// 但是事实上这个很少变，所以你可以延迟到必要的时候再注入
		bucket: ekit.ToPtr[string]("webook-1314583317"),
		daoArticle: daoArticle{
			db: db,
		},
	}
}

// 因为s
func (o *S3DAO) Sync(ctx context.Context, art Article) (int64, error) {
	// 保存制作库
	// 保存线上库，并且把 content 上传到 OSS
	// 制作库流量不大，并发不高，你就保存到数据库就可以
	// 当然，有钱或者体量大，就还是考虑 OSS
	var (
		id = art.Id
	)
	err := o.db.Transaction(func(tx *gorm.DB) error {
		var err error
		now := time.Now().UnixMilli()
		txDao := NewDaoArticle(tx)
		//制作库
		if id > 0 {
			err = txDao.Update(ctx, art)
		} else {
			id, err = txDao.Create(ctx, art)
		}
		if err != nil {
			return err
		}
		art.Id = id
		//线上库
		publishArt := PublishedArticleV1{
			Id:       art.Id,
			Title:    art.Title,
			AuthorId: art.AuthorId,
			Status:   art.Status,
			Ctime:    now,
			Utime:    now,
		}
		// 线上库不保存 Content,要准备上传到 OSS 里面
		return tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"title":  art.Title,
				"utime":  now,
				"status": art.Status,
			}),
		}).Create(&publishArt).Error
	})
	// 说明保存到数据库的时候失败了
	if err != nil {
		return 0, err
	}
	// 接下来就是保存到 OSS 里面
	// 你要有监控，你要有重试，你要有补偿机制
	_, err = o.oss.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:      o.bucket,
		Key:         ekit.ToPtr[string](strconv.FormatInt(art.Id, 10)),
		Body:        bytes.NewReader([]byte(art.Content)),
		ContentType: ekit.ToPtr[string]("text/plain;charset=utf-8"),
	})
	return id, err
}
func (o *S3DAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	//panic("implement me")
	err := o.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Model(&Article{}).Where("id = ? AND author_id = ?", id, author).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			errors.New("用户在尝试操作非本人数据")
		}
		res = tx.Model(&PublishedArticle{}).Where("id = ? AND author_id = ?", id, author).Update("status", status)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			errors.New("用户在尝试操作非本人数据")
		}
		return nil
	})
	if err != nil {
		return err
	}
	if status == statusPrivate {
		_, err = o.oss.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
			Bucket: o.bucket,
			Key:    ekit.ToPtr[string](strconv.FormatInt(id, 10)),
		})
	}
	return err
}
