// 生产者
package article

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
)

type ProducerEvents interface {
	ProduceReadEvent(ctx context.Context, e ReadEvent) error
	ProduceReadEventV1(ctx context.Context, v1 ReadEventv1)
}
type producerEvents struct {
	producer sarama.SyncProducer
}

func NewProducerEvents(producer sarama.SyncProducer) ProducerEvents {
	return &producerEvents{
		producer: producer,
	}
}
func (p *producerEvents) ProduceReadEvent(ctx context.Context, e ReadEvent) error {
	//给ReadEvent序列化
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	//发送消息
	_, _, err = p.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "read_article",
		Value: sarama.ByteEncoder(data),
	})
	return err
}
func (k *producerEvents) ProduceReadEventV1(ctx context.Context, v1 ReadEventv1) {
	//TODO implement me
	panic("implement me")
}

type ReadEvent struct {
	Uid int64 //用户id
	Aid int64 //文章id
}

// 批量
type ReadEventv1 struct {
	Uid []int64 //用户id
	Aid []int64
}
