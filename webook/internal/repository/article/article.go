package article

import (
	"context"
	"github.com/ecodeclub/ekit/slice"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository"
	"goFoundation/webook/internal/repository/cache"
	"goFoundation/webook/internal/repository/dao/article"
	"goFoundation/webook/pkg/logger"
	"gorm.io/gorm"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, id int64, author int64, status uint8) error
	List(ctx context.Context, id int64, Offset, Limit int) ([]domain.Article, error)
	Detail(ctx context.Context, id int64) (domain.Article, error)
	PubDetail(ctx context.Context, id int64) (domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error)
}
type articleRepository struct {
	dao     article.DaoArticle
	userDao repository.UserRepository
	author  article.DaoArticleAuthor
	reader  article.DaoArticleReader
	db      *gorm.DB
	cmd     cache.ArticleCache
	l       logger.LoggerV1
}

func NewArticleRepository(dao article.DaoArticle, userDao repository.UserRepository, cmd cache.ArticleCache, l logger.LoggerV1) ArticleRepository {
	return &articleRepository{
		dao:     dao,
		userDao: userDao,
		cmd:     cmd,
		l:       l,
	}
}
func NewArticleRepositoryV1(author article.DaoArticleAuthor, reader article.DaoArticleReader) ArticleRepository {
	return &articleRepository{
		author: author,
		reader: reader,
	}
}

// 热搜
func (a *articleRepository) ListPub(ctx context.Context, start time.Time, offset, limit int) ([]domain.Article, error) {
	res, err := a.dao.ListPub(ctx, start, offset, limit)
	if err != nil {
		return nil, err
	}
	return slice.Map(res, func(idx int, src article.Article) domain.Article {
		return a.toDomain(src)
	}), nil

}

// 创作者查id
func (a *articleRepository) Detail(ctx context.Context, id int64) (domain.Article, error) {
	res, err := a.cmd.Get(ctx, id)
	if err == nil {
		return res, err
	}
	art, err := a.dao.Detail(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	return a.toDomain(art), nil
}

// 读者id 在repo这做redis缓存
func (a *articleRepository) PubDetail(ctx context.Context, id int64) (domain.Article, error) {
	res, err := a.cmd.GetPub(ctx, id)
	if err == nil {
		return res, err
	}
	// 读取线上库数据，如果你的 Content 被你放过去了 OSS 上，你就要让前端去读 Content 字段
	art, err := a.dao.PubDetail(ctx, id)
	if err != nil {
		return domain.Article{}, err
	}
	//组装user 为了让读者看到作者的名字 适合单体应用
	usr, err := a.userDao.FindById(ctx, art.Id)
	res = domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Status:  domain.ArticleStatus(art.Status),
		Content: art.Content,
		Author: domain.Author{
			Id:   usr.Id,
			Name: usr.Nickname,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
	//做同步
	go func() {
		if err = a.cmd.SetPub(ctx, res); err != nil {
			a.l.Error("缓存已发表文章失败",
				logger.Error(err), logger.Int64("aid", res.Id))
		}
	}()
	return res, nil
}

// 搜素
func (a *articleRepository) List(ctx context.Context, id int64, Offset, Limit int) ([]domain.Article, error) {
	//先查缓存在查数据库
	if Limit == 0 && Offset <= 100 {
		res, err := a.cmd.GetFirstPage(ctx, id)
		if err == nil {
			go func() {
				//就是我在先查 不管查到多少页 我就先预加载缓存第一页 人的习惯就是只看第一页
				a.preCache(ctx, res)
			}()
			return res, nil
		}
		if err != cache.ErrKeyNotExist {
			a.l.Error("查询缓存文章失败",
				logger.Int64("author", id), logger.Error(err))
		}

	}
	res, err := a.dao.List(ctx, id, Offset, Limit)
	if err != nil {
		return nil, err
	}
	data := slice.Map[article.Article, domain.Article](res, func(idx int, src article.Article) domain.Article {
		return a.toDomain(src)
	})
	// 一般都是让调用者来控制是否异步。
	go func() {
		a.preCache(ctx, data)
	}()
	//回显 并发高用del 否则用set
	err = a.cmd.SetFirstPage(ctx, id, data)
	if err != nil {
		a.l.Error("刷新第一页文章的缓存失败",
			logger.Int64("author", id), logger.Error(err))
	}
	return data, nil
}

// 预加载实现 因为这是一个预测性质的，所以过期时间设置得很短
func (a *articleRepository) preCache(ctx context.Context, arts []domain.Article) {
	const contentSizeThreshold = 1024 * 1024
	if len(arts) > 0 && len(arts[0].Content) <= contentSizeThreshold {
		//打印日志
		err := a.cmd.Set(ctx, arts[0])
		if err != nil {
			a.l.Error("提前准备缓存失败", logger.Error(err))
		}
	}
}

// 仅自己看状态
func (a *articleRepository) Withdraw(ctx context.Context, id, author int64, status uint8) error {
	return a.dao.Withdraw(ctx, id, author, uint8(domain.ArticleStatusPrivate))
}

// 发表
func (a *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.dao.Sync(ctx, a.toEntity(art))
	if err != nil {
		return 0, err
	}
	go func() {
		author := art.Author.Id
		err = a.cmd.DelFirstPage(ctx, author)
		if err != nil {
			a.l.Error("删除第一页缓存失败",
				logger.Int64("author", author), logger.Error(err))
		}
	}()
	return id, nil
}

// 发表 Repository同步数据
func (a *articleRepository) SyncV1(ctx context.Context, art domain.Article) (int64, error) {
	var (
		id  = art.Id
		err error
	)
	if id > 0 {
		err = a.author.Update(ctx, a.toEntity(art))
	} else {
		id, err = a.author.Create(ctx, a.toEntity(art))
	}
	if err != nil {
		return 0, nil
	}
	art.Id = id
	err = a.reader.Save(ctx, a.toEntity(art))
	return id, err
}

// Repository同步数据 gorm事务 这种不建议使用
func (a *articleRepository) SyncV2(ctx context.Context, art domain.Article) (int64, error) {
	tx := a.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return 0, tx.Error
	}
	// 直接 defer Rollback
	// 如果我们后续 Commit 了，这里会得到一个错误，但是没关系
	defer tx.Rollback()
	authorDao := article.NewDaoArticleAuthor(tx)
	readerDao := article.NewDaoArticleReader(tx)
	var (
		id  = art.Id
		err error
	)
	artn := a.toEntity(art)
	if id > 0 {
		err = authorDao.Update(ctx, artn)
	} else {
		id, err = authorDao.Create(ctx, artn)
	}
	if err != nil {
		return 0, nil
	}
	art.Id = id
	err = readerDao.SaveV2(ctx, article.PublishedArticle(artn))
	if err != nil {
		// 依赖于 defer 来 rollback
		return 0, err
	}
	tx.Commit()
	return art.Id, err
}

// 创建
func (a *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	id, err := a.dao.Create(ctx, a.toEntity(art))
	if err != nil {
		return 0, err
	}
	//删除缓存
	author := art.Author.Id
	go func() {
		err = a.cmd.DelFirstPage(ctx, author)
		if err != nil {
			a.l.Error("删除缓存失败",
				logger.Int64("author", author), logger.Error(err))
		}
	}()
	return id, nil
}

// 修改
func (a *articleRepository) Update(ctx context.Context, art domain.Article) error {
	err := a.dao.Update(ctx, a.toEntity(art))
	if err != nil {
		return err
	}
	//删除缓存
	author := art.Author.Id
	go func() {
		err = a.cmd.DelFirstPage(ctx, author)
		if err != nil {
			a.l.Error("删除缓存失败",
				logger.Int64("author", author), logger.Error(err))
		}
	}()
	return nil
}

// 转换类型
func (a *articleRepository) toEntity(art domain.Article) article.Article {
	return article.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
		Status:   uint8(art.Status),
	}
}

// 转换类型
func (a *articleRepository) toDomain(art article.Article) domain.Article {
	return domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		Ctime: time.UnixMilli(art.Ctime),
		Utime: time.UnixMilli(art.Utime),
	}
}

//我刚接收我们公司的系统的时候 出过一个故障 这个故障是大对象的故障 redis内存不够 因为我们有些业务他的正常的数据都不大
//但是在遇到一些大客户 大v的是时候数据量特别大 很快就把redis占满 在这种情况下 我们控制redis缓存策略也就是在回显的redis的时候
//当某一个值 某个用户的数据它的大小超过一定的范围 我就不会执行这个缓存
//我把宝贵的redis内存腾出来给那些小的对象去用 不是说所有情况下都不要缓存大对象，
//而是说，你要权衡性能和内存开销 如果大对象 他的计算慢查询慢也是要缓存大对象的
