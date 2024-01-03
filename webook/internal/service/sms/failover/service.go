// 自动切换不同服务商
package failover

import (
	"context"
	"errors"
	"goFoundation/webook/internal/service/sms"
	"log"
	"sync/atomic"
)

type FailoverSmsService struct {
	//有很多厂商
	svcs []sms.Service
	idx  uint64
}

func NewFailoverSmsService(svcs []sms.Service) sms.Service {
	return &FailoverSmsService{
		svcs: svcs,
	}
}

// 轮询切换
func (f *FailoverSmsService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//这种方式 你每次轮询完事都从第一个开始
	for _, v := range f.svcs {
		err := v.Send(ctx, tpl, args, numbers...)
		if err == nil {
			return nil
		}
		//正常这边，输出日志
		//做监控
		log.Println(err)
	}
	return errors.New("所以厂商都失败了")
}

// 严格轮询
func (f *FailoverSmsService) SendV1(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//取下一个节点作为起始节点
	idx := atomic.AddUint64(&f.idx, 1)
	length := uint64(len(f.svcs))
	for i := idx; i < idx+length; i++ {
		//不能直接f.svcs【i】 是因为f.svcs+idx大于了这本身长度
		svc := f.svcs[int(i%length)]
		err := svc.Send(ctx, tpl, args, numbers...)
		switch err {
		case nil:
			return nil
		case context.DeadlineExceeded, context.Canceled: //超时
			return err
		default:
			// 输出日志
		}
	}
	return errors.New("所以厂商都失败了")
}
