package ioc

import (
	"github.com/IBM/sarama"
	"github.com/spf13/viper"
	events2 "goFoundation/webook/interactive/events"
	"goFoundation/webook/pkg/saramax"
)

func InitKafka() sarama.Client {
	type Config struct {
		Addrs []string `json:"addrs"`
	}
	//第一步定义
	s := sarama.NewConfig()
	//同步
	s.Producer.Return.Successes = true
	//反序列化
	var c Config
	err := viper.UnmarshalKey("kafka", &c)
	if err != nil {
		panic(err)
	}
	//
	client, err := sarama.NewClient(c.Addrs, s)
	if err != nil {
		panic(err)
	}
	return client
}
func NewSyncProducer(c sarama.Client) sarama.SyncProducer {
	res, err := sarama.NewSyncProducerFromClient(c)
	if err != nil {
		panic(err)
	}
	return res
}

// NewConsumers 面临的问题依旧是所有的 Consumer 在这里注册一下
//func NewConsumers(c *article.ConsumerEvent) []events.ConsumerEvents {
//	return []events.ConsumerEvents{c}
//}

// 批量
func NewConsumers(c *events2.InteractiveReadEventBatchConsumer) []saramax.ConsumerEvents {
	return []saramax.ConsumerEvents{c}
}
