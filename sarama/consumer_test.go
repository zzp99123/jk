// 消费者
package sarama

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
	"log"
	"testing"
	"time"
)

func TestConsumer(t *testing.T) {
	cfg := sarama.NewConfig()
	// 正常来说，一个消费者都是归属于一个消费者的组的
	// 消费者组就是你的业务
	consumer, err := sarama.NewConsumerGroup(addrs,
		"test_group", cfg)
	assert.NoError(t, err)
	//带超时的context
	start := time.Now()
	//过期时间
	//ctx,cancel =context.WithTimeout(context.Background(),time.Minute * 10)
	//defer cancel()
	ctx, cancel := context.WithCancel(context.Background())
	time.AfterFunc(time.Second*5, func() {
		cancel()
	})
	err = consumer.Consume(ctx, []string{"test_topic"}, testConsumerGroupHandler{})
	// 你消费结束，就会到这里
	t.Log(err, time.Since(start).String())
}

type testConsumerGroupHandler struct{}

// 不能用指针
func (t testConsumerGroupHandler) Setup(s sarama.ConsumerGroupSession) error {
	//指定偏移量 在部分场景下，我们会希望消费历史消息，或者从某个消息开始消费，那么可以考虑在 Setup 里面设置偏移量。
	//不过一般我都是建议走离线渠道，操作 Kafka 集群去重置对应的偏移量。 核心在于，你并不是每次重新部署，重新启动都是要重置这个偏移量的。
	//关键调用是 ResetOffset
	// topic => 偏移量
	topics := s.Claims()["test_topic"]
	for _, part := range topics {
		s.ResetOffset("test_topic", part, sarama.OffsetOldest, "") //偏移量-2
		//s.ResetOffset("test_topic", part, 123, "") 具体值
		//s.ResetOffset("test_topic", part, sarama.OffsetNewest, "") 偏移量-1
	}
	return nil
}
func (t testConsumerGroupHandler) Cleanup(s sarama.ConsumerGroupSession) error {
	log.Println("Cleanup")
	return nil
}

// 异步
func (t testConsumerGroupHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	//代表的是你和kafka的绘画 从建立链接到彻底断掉的那一段世时间
	msgs := c.Messages()
	//定义批次
	const batchSize = 10
	for {
		//定义超时时间
		ctx, cacel := context.WithTimeout(context.Background(), time.Second*5)
		var eg errgroup.Group
		//定义倒数第一次数据
		var last *sarama.ConsumerMessage
		done := false
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				//超时了
				done = true
				cacel()
				return nil
			case msg, ok := <-msgs:
				if !ok {
					cacel()
					// 代表消费者被关闭了
					return nil
				}
				last = msg
				eg.Go(func() error {
					// 我就在这里消费
					time.Sleep(time.Second)
					// 你在这里重试
					log.Println(string(msg.Value))
					return nil
				})

			}
			if done {
				break
			}
		}
		cacel()
		err := eg.Wait()
		if err != nil {
			// 你在这里重试
			//打日志
			continue
		}
		if last != nil {
			s.MarkMessage(last, "") //不用提交全部 提交最后一个就行
		}

	}

}

// 同步
func (t testConsumerGroupHandler) ConsumeClaimv1(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	//代表的是你和kafka的绘画 从建立链接到彻底断掉的那一段世时间
	msgs := c.Messages()
	for m := range msgs {
		//var my MyBizMsg
		//err := json.Unmarshal(m.Value, &my)
		//if err != nil {
		//	//这就是消费消息出错
		//	//重试
		//	//打日志
		//	continue
		//}
		log.Println(string(m.Value))
		s.MarkMessage(m, "")
	}
	return nil
}

type MyBizMsg struct {
	Name string
}

// 返回只读channel 这个用的多
func chanelv1() <-chan struct{} {
	panic("implement me")
}

// 返回读写channel 偶尔
func chanelv2() chan struct{} {
	panic("implement me")
}

// 返回只写channel 很少用
func chanelv3() chan<- struct{} {
	panic("implement me")
}
