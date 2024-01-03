package article

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type mongoDBDAO struct {
	//client *mongo.Client
	// 代表 webook 的
	//database *mongo.Database
	// 代表的是制作库
	col *mongo.Collection
	// 代表的是线上库
	liveCol *mongo.Collection
	node    *snowflake.Node
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) DaoArticle {
	return &mongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}
func InitCollections(db *mongo.Database) error {
	//定义过期时间
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().
		CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").Indexes().
		CreateMany(ctx, index)
	return err
}
func (m *mongoDBDAO) Create(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	//id := m.idGen()
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	// 你没有自增主键
	// GLOBAL UNIFY ID (GUID，全局唯一ID）
	return id, err
}
func (m *mongoDBDAO) Update(ctx context.Context, art Article) error {
	//相当于sql的whrer
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.D{bson.E{"$set", bson.M{
		"title":   art.Title,
		"content": art.Content,
		"utime":   time.Now().UnixMilli(),
		"status":  art.Status,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	// 这边就是校验了 author_id 是不是正确的 ID
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}

// 发表
func (m *mongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = m.Update(ctx, art)
	} else {
		id, err = m.Create(ctx, art)
	}
	if err != nil {
		return id, err
	}
	art.Id = id
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	now := time.Now().UnixMilli()
	art.Utime = now
	_, err = m.liveCol.UpdateOne(ctx, filter, bson.D{bson.E{"$set", art},
		bson.E{"$setOnInsert", bson.D{bson.E{"ctime", now}}}},
		options.Update().SetUpsert(true),
	)
	return id, err
}

// 仅自己可见
func (m *mongoDBDAO) Withdraw(ctx context.Context, id int64, author int64, status uint8) error {
	//panic("implement me")
	filter := bson.M{"id": id, "author_id": author}
	sets := bson.D{bson.E{"status", status}}
	res, err := m.col.UpdateOne(ctx, filter, bson.D{bson.E{"$set", sets}})
	if err != nil {
		return err
	}
	// 这边就是校验了 author_id 是不是正确的 ID
	if res.ModifiedCount == 0 {
		return errors.New("更新数据失败")
	}
	return nil
}
func (m *mongoDBDAO) List(ctx context.Context, id int64, Offset, Limit int) ([]Article, error) {
	panic("import me")
}
func (m *mongoDBDAO) Detail(ctx context.Context, id int64) (Article, error) {
	panic("import me")
}
func (m *mongoDBDAO) PubDetail(ctx context.Context, id int64) (PublishedArticle, error) {
	panic("import me")
}
func (m *mongoDBDAO) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]Article, error) {
	panic("import me")
}
