// 热搜
package service

import (
	"context"
	"errors"
	"github.com/ecodeclub/ekit/queue"
	"github.com/ecodeclub/ekit/slice"
	intrv1 "goFoundation/webook/api/proto/gen/intr/v1"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository"
	"math"
	"time"
)

type RangingService interface {
	TopN(ctx context.Context) error
	topN(ctx context.Context) ([]domain.Article, error)
}
type rangingService struct {
	artSvc    ArticleService
	intrSvc   intrv1.InteractiveServiceClient
	repo      repository.RangingRepository
	batchSize int //Limit 指定要查询的最大记录数
	n         int
	// scoreFunc 不能返回负数
	scoreFunc func(t time.Time, likeCnt int64) float64 //计算热度的算法
	// 负载
	load int64
}

func NewRangingService(artSvc ArticleService, intrSvc intrv1.InteractiveServiceClient) RangingService {
	return &rangingService{
		artSvc:    artSvc,
		intrSvc:   intrSvc,
		batchSize: 100,
		n:         100,
		scoreFunc: func(t time.Time, likeCnt int64) float64 {
			return float64(likeCnt-1) / math.Pow(float64(likeCnt+2), 1.5)
		},
	}
}
func (r *rangingService) TopN(ctx context.Context) error {
	res, err := r.topN(ctx)
	if err != nil {
		return err
	}
	// 在这里，存起来 利用redis缓存
	return r.repo.TopN(ctx, res)
}

// 方便测试用返回值([]domain.Article, error)
func (r *rangingService) topN(ctx context.Context) ([]domain.Article, error) {
	// 我只取七天内的数据
	now := time.Now()
	type Score struct {
		art   domain.Article
		score float64
	}
	// 这里可以用非并发安全
	topN := queue.NewConcurrentPriorityQueue[Score](r.n, func(src Score, dst Score) int {
		if src.score > dst.score {
			return 1
		} else if src.score == dst.score {
			return 0
		} else {
			return -1
		}
	})
	offset := 0 //Offset指定开始返回记录前要跳过的记录数
	for {
		//从数据库中拉取一批文章
		res, err := r.artSvc.ListPub(ctx, now, offset, r.batchSize)
		if err != nil {
			return nil, err
		}
		//我要从这批文章里把id提取出来
		idx := slice.Map[domain.Article, int64](res, func(idx int, src domain.Article) int64 {
			return src.Id
		})
		//查id并且要去找到对应的点赞数据
		ids, err := r.intrSvc.GetByIds(ctx, &intrv1.GetByIdsRequest{
			Biz: "article", BizIds: idx,
		})
		if err != nil {
			return nil, err
		}
		if len(ids.Intr) == 0 {
			return nil, errors.New("没有数据")
		}
		//排序
		for _, v := range res {
			//根据每个文章的id找到相对应的点赞数
			id := ids.Intr[v.Id]
			s := r.scoreFunc(v.Utime, id.LikeCnt)
			// 我要考虑，我这个 score 在不在前一百名
			// 拿到热度最低的
			err = topN.Enqueue(Score{
				art:   v,
				score: s,
			})
			// 这种写法，要求 topN 已经满了
			//if err == queue.ErrOutOfCapacity {
			//	val, _ := topN.Dequeue()
			//	if val.score < s {
			//		err = topN.Enqueue(Score{
			//			art:   v,
			//			score: s,
			//		})
			//	} else {
			//		_ = topN.Enqueue(val)
			//	}
			//}
		}
		// 一批已经处理完了，问题来了，我要不要进入下一批？我怎么知道还有没有？
		//一天一更新热搜就跟微博热搜一样
		if len(res) < r.batchSize || now.Sub(res[len(res)-1].Utime).Hours() > 24 {
			// 我这一批都没取够，我当然可以肯定没有下一批了
			break
		}
		//更新Offset
		offset = offset + len(res)
	}
	// 最后得出结果
	val := make([]domain.Article, r.n)

	for i := r.n - 1; i >= 0; i-- {
		vals, err := topN.Dequeue()
		if err != nil {
			// 说明取完了，不够 n
			break
		}
		val[i] = vals.art
	}
	return val, nil
}
