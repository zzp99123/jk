// 批量消费
package events

import (
	"context"
	"github.com/IBM/sarama"
	"goFoundation/webook/interactive/repository"
	"goFoundation/webook/pkg/logger"
	"goFoundation/webook/pkg/saramax"
	"time"
)

type InteractiveReadEventBatchConsumer struct {
	client sarama.Client
	repo   repository.InteractiveRepository
	l      logger.LoggerV1
}

func NewInteractiveReadEventBatchConsumer(client sarama.Client, l logger.LoggerV1, repo repository.InteractiveRepository) *InteractiveReadEventBatchConsumer {
	return &InteractiveReadEventBatchConsumer{
		client: client,
		l:      l,
		repo:   repo,
	}
}
func (i *InteractiveReadEventBatchConsumer) Start() error {
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		i.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewBatchHandler[ReadEvent](i.l, i.Consume))
		if err != nil {
			i.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (i *InteractiveReadEventBatchConsumer) Consume(msg []*sarama.ConsumerMessage, ts []ReadEvent) error {
	ids := make([]int64, 0, len(ts))
	bizs := make([]string, 0, len(ts))
	for _, evt := range ts {
		ids = append(ids, evt.Aid)
		bizs = append(bizs, "article")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	err := i.repo.BatchIncrReadCnt(ctx, bizs, ids)
	if err != nil {
		i.l.Error("批量增加阅读计数失败",
			logger.Field{Key: "ids", Value: ids},
			logger.Error(err))
	}
	return nil

}

//我参与的项目有个严重的问题 就是生产者/消费者性能比较差 产生消息积压 我当时解决这个问题2个思路 一个是异步消费 一个是批量接口
//第一个 异步消费开groution 为了解决消息偏移量的问题 我们搞的是 异步消费 批量提交 last组织一下细节 我们把消息积压解决好了以后 数据库受不了然后进一步提出另一个方案批量接口 因为我们消费者这边可以把单个改成批量 于是我们中和了异步和消费者批量接口优化了性能同时数据库的压力也不大

//难点 消息积压  我需要异步消费 消息消费了但没提交怎么办 批量消费 如果批量提交了都没提交上怎么办 重试 幂等 我可以利用redis高可用 用redis来去重 唯一索引
//虽热我解决了 消息积压问题 但是数据库压力大了导致查询变慢了影响到其他业务了然后 批量接口 为什么批量就快了 从mvcc的原理来回答 vdolog andolog
