package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"goFoundation/webook/pkg/logger"
	"log"
)

type Handler[T any] struct {
	l  logger.LoggerV1
	fn func(msg *sarama.ConsumerMessage, r T) error
	//vec
}

func NewHandler[T any](l logger.LoggerV1, fn func(msg *sarama.ConsumerMessage, r T) error) *Handler[T] {
	return &Handler[T]{
		l:  l,
		fn: fn,
	}
}

func (h *Handler[T]) Setup(s sarama.ConsumerGroupSession) error {
	log.Println("Setup")
	return nil
}
func (h *Handler[T]) Cleanup(s sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}
func (h Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgs := claim.Messages()
	for msg := range msgs {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			h.l.Error("反序列化消息失败",
				logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int64("partition", int64(msg.Partition)),
				logger.Int64("offset", msg.Offset))
			continue
		}
		// 在这里执行重试
		for i := 0; i < 3; i++ {
			err = h.fn(msg, t)
			if err == nil {
				break
			}
			h.l.Error("处理消息失败",
				logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int64("partition", int64(msg.Partition)),
				logger.Int64("offset", msg.Offset))
		}

		if err != nil {
			h.l.Error("处理消息失败-重试次数上限",
				logger.Error(err),
				logger.String("topic", msg.Topic),
				logger.Int64("partition", int64(msg.Partition)),
				logger.Int64("offset", msg.Offset))
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
