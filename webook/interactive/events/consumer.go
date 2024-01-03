// 消费者
package events

import (
	"context"
	"github.com/IBM/sarama"
	"goFoundation/webook/interactive/repository"
	"goFoundation/webook/pkg/logger"
	"goFoundation/webook/pkg/saramax"
	"time"
)

type ConsumerEvent struct {
	client sarama.Client //集合Client接口
	l      logger.LoggerV1
	repo   repository.InteractiveRepository
}

func NewConsumerEvents(client sarama.Client, l logger.LoggerV1, repo repository.InteractiveRepository) *ConsumerEvent {
	return &ConsumerEvent{
		client: client,
		l:      l,
		repo:   repo,
	}
}

func (c *ConsumerEvent) Start() error {
	//NewConsumerGroupFromClient这个不用填参数adress
	cg, err := sarama.NewConsumerGroupFromClient("interactive",
		c.client)
	if err != nil {
		return err
	}
	go func() {
		err := cg.Consume(context.Background(),
			[]string{"read_article"},
			saramax.NewHandler[ReadEvent](c.l, c.Consume))
		if err != nil {
			c.l.Error("退出了消费循环异常", logger.Error(err))
		}
	}()
	return err
}
func (c *ConsumerEvent) Consume(msg *sarama.ConsumerMessage, r ReadEvent) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return c.repo.IncrReadCnt(ctx, "article", r.Aid)
}
