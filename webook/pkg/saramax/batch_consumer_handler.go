// 二次封装sarama 批量发送
package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"goFoundation/webook/pkg/logger"
	"log"
	"time"
)

type BatchHandler[T any] struct {
	l             logger.LoggerV1
	fn            func(msg []*sarama.ConsumerMessage, t []T) error //批次是多个 所以用切片
	batchSize     int                                              //批次
	batchDuration time.Duration                                    //过期时间
}

func NewBatchHandler[T any](l logger.LoggerV1, fn func(msg []*sarama.ConsumerMessage, t []T) error) *BatchHandler[T] {
	return &BatchHandler[T]{
		l:             l,
		fn:            fn,
		batchDuration: time.Second,
		batchSize:     10,
	}
}

func (b BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	log.Println("Setup")
	return nil
}

func (b BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}

func (b BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgch := claim.Messages()
	for {
		//cancel 就可以告诉 ctx，这个 goroutine 应该退出了，从而进入 case <-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		done := false                                           //控制这个cancel()
		ba := make([]T, 0, b.batchSize)                         //一个批次放100条
		msgs := make([]*sarama.ConsumerMessage, 0, b.batchSize) //把n个批次都进[]*sarama.ConsumerMessage
		for i := 0; i < b.batchSize && !done; i++ {
			select {
			case <-ctx.Done(): //用来在 goroutine 内部用来判断当前 goroutine 是否应该退出。
				done = true
			case msg, ok := <-msgch:
				if !ok {
					cancel()
					// 代表消费者被关闭了
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {
					b.l.Error("反序列化失败",
						logger.Error(err),
						logger.String("topic", msg.Topic),
						logger.Int64("partition", int64(msg.Partition)),
						logger.Int64("offset", msg.Offset))
					continue

				}
				msgs = append(msgs, msg) //用来将元素添加到切片末尾并返回结果。
				ba = append(ba, t)
			}
		}
		cancel()
		if len(msgs) == 0 {
			//没有数据
			continue
		}
		err := b.fn(msgs, ba)
		if err != nil {
			b.l.Error("调用业务批量接口失败",
				logger.Error(err))
			// 你这里整个批次都要记下来
			// 还要继续往前消费
		}
		for _, msg := range msgs {
			// 处理消息成功后标记为处理, 然后会自动提交
			session.MarkMessage(msg, "")
		}
	}

}
