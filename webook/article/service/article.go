package service

import (
	"context"
	"goFoundation/webook/internal/domain"
	events "goFoundation/webook/internal/events/article"
	"goFoundation/webook/internal/repository/article"
	"goFoundation/webook/pkg/logger"
	"time"
)

//go:generate mockgen -source=article.go -package=svcmocks -destination=mocks/article.mock.go ArticleService
type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, art domain.Article) error
	List(ctx context.Context, id int64, Offset, Limit int) ([]domain.Article, error)
	Detail(ctx context.Context, id int64) (domain.Article, error)
	PubDetail(ctx context.Context, id, uid int64) (domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error)
}

type articleService struct {
	// 2. 在 repo 里面处理制作库和线上库
	// 1 和 2 是互斥的，不会同时存在
	repo article.ArticleRepository

	author   article.RepositoryArticleAuthor
	reader   article.RepositoryArticleReader
	l        logger.LoggerV1
	producer events.ProducerEvents
	ch       chan readChan
}
type readChan struct {
	aid int64
	uid int64
}

func NewArticleService(repo article.ArticleRepository, l logger.LoggerV1, producer events.ProducerEvents) ArticleService {
	return &articleService{
		repo:     repo,
		l:        l,
		producer: producer,
	}
}
func NewArticleServiceV1(repo article.ArticleRepository,
	l logger.LoggerV1, author article.RepositoryArticleAuthor, reader article.RepositoryArticleReader) ArticleService {
	return &articleService{
		repo:   repo,
		l:      l,
		author: author,
		reader: reader,
	}
}

// 批量生产者 提高性能
func NewArticleServiceV2(repo article.ArticleRepository,
	l logger.LoggerV1, producer events.ProducerEvents) ArticleService {
	ch := make(chan readChan, 10)
	go func() {
		for {
			aids := make([]int64, 0, 10)
			uids := make([]int64, 0, 10)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			for i := 0; i < 10; i++ {
				select {
				case infor, ok := <-ch:
					if !ok {
						cancel()
						return
					}
					aids = append(aids, infor.aid)
					uids = append(uids, infor.uid)
				case <-ctx.Done():
					break
				}
			}
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			producer.ProduceReadEventV1(ctx, events.ReadEventv1{
				Uid: uids,
				Aid: aids,
			})
			cancel()
		}
	}()
	return &articleService{
		repo:     repo,
		l:        l,
		ch:       ch,
		producer: producer,
	}
}

// 热搜
func (a *articleService) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error) {
	return a.repo.ListPub(ctx, start, offset, limit)
}

// 创作者查id
func (a *articleService) Detail(ctx context.Context, id int64) (domain.Article, error) {
	return a.repo.Detail(ctx, id)
}
func (a *articleService) PubDetail(ctx context.Context, id, uid int64) (domain.Article, error) {
	art, err := a.repo.PubDetail(ctx, id)
	//利用kafka 让读者读过的文章标记下来
	//也可以批量消费者 批量消费者不是单纯的for循环 而是利用管道把uid,aid做成切片 传输进去
	if err == nil {
		go func() {
			er := a.producer.ProduceReadEvent(ctx, events.ReadEvent{
				Aid: id,
				Uid: uid,
			})
			if er == nil {
				a.l.Error("发送读者阅读事件失败")
			}
		}()
	}
	return art, err
}

// 搜素
func (a *articleService) List(ctx context.Context, id int64, Offset, Limit int) ([]domain.Article, error) {
	return a.repo.List(ctx, id, Offset, Limit)
}

// 仅自己查看状态
func (a *articleService) Withdraw(ctx context.Context, art domain.Article) error {
	//art.Status = domain.ArticleStatusUnpublished
	return a.repo.Withdraw(ctx, art.Id, art.Author.Id, uint8(domain.ArticleStatusPrivate))
}
func (a *articleService) Save(ctx context.Context,
	art domain.Article) (int64, error) {
	//定义未发表状态
	art.Status = domain.ArticleStatusUnpublished
	if art.Id > 0 {
		err := a.update(ctx, art)
		return art.Id, err
	}
	return a.create(ctx, art)
}

// 发表
func (a *articleService) Publish(ctx context.Context,
	art domain.Article) (int64, error) {
	//定义发表状态
	art.Status = domain.ArticleStatusPublished
	return a.repo.Sync(ctx, art)
}

// 发表 Service 层同步数据写法
func (a *articleService) PublishV1(ctx context.Context,
	art domain.Article) (int64, error) {
	//确保制作库与线上库的id相等
	var (
		id  = art.Id
		err error
	)
	//art.Id大于0就是修改 否则就是创建制作库 然后
	if art.Id > 0 {
		err = a.author.Update(ctx, art)
		//return art.Id, err
	} else {
		id, err = a.author.Create(ctx, art)
	}
	if err != nil {
		return 0, err
	}
	//确保制作库与线上库的id相等
	art.Id = id
	//重试
	for i := 0; i < 3; i++ {
		err = a.reader.Save(ctx, art)
		if err == nil {
			break
		}
		a.l.Error("部分失败，保存到线上库失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
	}
	if err != nil {
		a.l.Error("部分失败，重试彻底失败",
			logger.Int64("art_id", art.Id),
			logger.Error(err))
		// 接入你的告警系统，手工处理一下
		// 走异步，我直接保存到本地文件
		// 走 Canal
		// 打 MQ
	}
	return id, err
}

// 创建
func (a *articleService) create(ctx context.Context,
	art domain.Article) (int64, error) {
	return a.repo.Create(ctx, art)
}

// 修改
func (a *articleService) update(ctx context.Context,
	art domain.Article) error {
	return a.repo.Update(ctx, art)
}
