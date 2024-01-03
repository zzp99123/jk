package mongo

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestMongo(t *testing.T) {
	//初始化客户端
	//定义一个过期时间
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	monitor := &event.CommandMonitor{
		// 每个命令（查询）执行之前
		Started: func(ctx context.Context, startedEvent *event.CommandStartedEvent) {
			fmt.Println(startedEvent.Command)
		},
		//成功命令
		Succeeded: func(ctx context.Context, succeededEvent *event.CommandSucceededEvent) {

		},
		//失败命令
		Failed: func(ctx context.Context, failedEvent *event.CommandFailedEvent) {

		},
	}
	opts := options.Client().ApplyURI("mongodb://root:example@localhost:27017").SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opts)
	assert.NoError(t, err)
	//获取一个Databaes
	mdb := client.Database("webook")
	//在Databaes中获取一个集合
	col := mdb.Collection("articles")
	defer func() {
		_, err = col.DeleteMany(ctx, bson.D{})
	}()

	//插入
	res, err := col.InsertOne(ctx, Article{
		Id:      123,
		Title:   "我的标题",
		Content: "我的内容",
	})
	assert.NoError(t, err)
	// 这个是文档ID，也就是 mongodb 中的 _id 字段
	fmt.Printf("id %s", res.InsertedID)

	//查找
	filter := bson.D{bson.E{Key: "id", Value: 123}}
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	assert.NoError(t, err)
	fmt.Printf("%#v \n", art)

	art = Article{}
	err = col.FindOne(ctx, Article{
		Id: 123,
	}).Decode(&art)
	if err == mongo.ErrNoDocuments {
		fmt.Println("没有数据")
	}
	assert.NoError(t, err)
	fmt.Printf("%#v \n", art)

	//复杂查找 or
	os := bson.D{bson.E{"$or", bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"id", 456}}}}}
	f, err := col.Find(ctx, os)
	assert.NoError(t, err)
	var a []Article
	err = f.All(ctx, &a)
	assert.NoError(t, err)

	//复杂查找 and
	aos := bson.D{bson.E{"$and", bson.A{bson.D{bson.E{"id", 123}},
		bson.D{bson.E{"title", "我的标题"}}}}}
	and, err := col.Find(ctx, aos)
	assert.NoError(t, err)
	a = []Article{}
	err = and.All(ctx, &a)
	assert.NoError(t, err)

	//复杂查找 in
	o := bson.D{bson.E{"id", bson.D{bson.E{"$in", []int{123, 456}}}}}
	inres, err := col.Find(ctx, o)
	assert.NoError(t, err)
	var ain []Article
	err = inres.All(ctx, &ain)
	assert.NoError(t, err)

	//查询特定字段
	ins, err := col.Find(ctx, o, options.Find().SetProjection(bson.M{
		"id":    1,
		"title": 1,
	}))
	ain = []Article{}
	err = ins.All(ctx, &ain)
	assert.NoError(t, err)

	//创建索引
	ires, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.M{"author_id": 1},
		},
	})
	assert.NoError(t, err)
	fmt.Println(ires)

	//更新
	sets := bson.D{bson.E{"$set", bson.E{"title", "标题"}}}
	u, err := col.UpdateOne(ctx, filter, sets) //更新一个
	if err != nil {
		panic(err)
	}
	fmt.Println("affected", u.ModifiedCount)

	//更新多个
	u, err = col.UpdateMany(ctx, filter, bson.D{
		bson.E{"$set", Article{
			Id:    123,
			Title: "ok",
		}}})

	//删除
	d, err := col.DeleteMany(ctx, filter)
	if err != nil {
		panic(err)
	}
	fmt.Println("affected", d.DeletedCount)
}

type Article struct {
	Id       int64  `bson:"id,omitempty"`
	Title    string `bson:"title,omitempty"`
	Content  string `bson:"content,omitempty"`
	AuthorId int64  `bson:"author_id,omitempty"`
	Status   uint8  `bson:"status,omitempty"`
	Ctime    int64  `bson:"ctime,omitempty"`
	Utime    int64  `bson:"utime,omitempty"`
}
