// 生产者
package sarama

import (
	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"testing"
)

var addrs = []string{"localhost:9094"}

// 同步
func TestSyncProducer(t *testing.T) {
	cfg := sarama.NewConfig()
	//同步的product一定要设置
	cfg.Producer.Return.Successes = true
	//指定分区
	//Hash 的参数，本质上是一个 Hash 的实现
	//cfg.Producer.Partitioner = sarama.NewHashPartitioner              //Hash：根据 key 的哈希值来筛选一个。
	//cfg.Producer.Partitioner = sarama.NewRandomPartitioner            //Random：随机挑一个。
	//cfg.Producer.Partitioner = sarama.NewRoundRobinPartitioner        //RoundRobin：轮询。
	//cfg.Producer.Partitioner = sarama.NewManualPartitioner            // 根据 Message 中的 partition 字段来选择。
	//cfg.Producer.Partitioner = sarama.NewConsistentCRCHashPartitioner // ConsistentCRC：一致性哈希，用的是 CRC16 算法。
	//cfg.Producer.Partitioner = sarama.NewCustomPartitioner()          // ：实际上不 Custom，而是自定义一部分
	//cfg.Producer.Partitioner = sarama.NewCustomHashPartitioner(func() hash.Hash32 {
	//
	//})
	pro, err := sarama.NewSyncProducer(addrs, cfg)
	assert.NoError(t, err)
	for i := 0; i < 100; i++ {
		_, _, err = pro.SendMessage(&sarama.ProducerMessage{
			Topic: "read_article",
			//消息数据本体转josn
			Value: sarama.StringEncoder(`{"aid": 1, "uid": 123}`),
			//会在生产者和消费者之间传递
			//Headers: []sarama.RecordHeader{
			//	{
			//		Key:   []byte("trace_id"),
			//		Value: []byte("123456"),
			//	},
			//},
			////只作用于发送过程
			//Metadata: "这是metadata",
		})
		assert.NoError(t, err)
	}
}

// 异步
//func TestAsyncProducer(t *testing.T) {
//	cfg := sarama.NewConfig()
//	cfg.Producer.Return.Successes = true
//	cfg.Producer.Return.Errors = true
//	pro, err := sarama.NewAsyncProducer(addrs, cfg)
//	assert.NoError(t, err)
//	msgcg := pro.Input()
//	msgcg <- &sarama.ProducerMessage{
//		Topic: "test_topic",
//		Key:   sarama.StringEncoder("old-123"),
//		Value: sarama.StringEncoder("Hello, 这是一条消息 A"),
//		Headers: []sarama.RecordHeader{{
//			Key:   []byte("trace_id"),
//			Value: []byte("123456"),
//		}},
//		Metadata: "这是metadata",
//	}
//	// 在实践中，一般是开另外一个 goroutine 来处理结果的
//	errch := pro.Errors()
//	such := pro.Successes()
//	// 如果两个情况都没发生，就会阻塞
//	select {
//	case err := <-errch:
//		val, _ := err.Msg.Value.Encode()
//		t.Log(err.Err, string(val))
//	case msg := <-such:
//		val, _ := msg.Value.Encode()
//		t.Log("成功了", string(val))
//	}
//}
