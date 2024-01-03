// 动态判定服务商状态
// 错误率：例如连续 N 个超时响应，错误率超过 10%。
// 响应时间增长率：例如响应时间从 100ms 突然变成 1s
package failover

import (
	"context"
	"goFoundation/webook/internal/service/sms"
	"sync/atomic"
)

type TimeoutFailoverSMSService struct {
	svcs []sms.Service
	//记录下标
	idx int32
	//超时的个数
	cnt int32
	//超时几个就切换
	threshold int32
}

func NewTimeoutFailoverSMSService() sms.Service {
	return &TimeoutFailoverSMSService{}
}

// 只要连续超过 N 个请求 超时了，就直接切换。
func (t *TimeoutFailoverSMSService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt32(&t.idx)
	cnt := atomic.LoadInt32(&t.cnt)
	//首页一进来 我就看他超时的个数 大与 超时几个就切换
	if cnt > t.threshold {
		//切换下标到下一个 //防止溢出 % int32(len(t.svcs))
		newIdx := (idx + 1) % int32(len(t.svcs))
		if atomic.CompareAndSwapInt32(&t.idx, idx, newIdx) {
			//成功往后挪了一位 超时个数就从0开始计算
			atomic.StoreInt32(&t.cnt, 0)
		}
		//else 这要写的并发问题 同时2个人另一个换成了
		//这相当于 我把下标换成了新的下标
		//idx = newIdx
		idx = atomic.LoadInt32(&t.idx)
	}
	svc := t.svcs[idx]
	err := svc.Send(ctx, tpl, args, numbers...)
	switch err {
	case context.DeadlineExceeded:
		atomic.StoreInt32(&t.cnt, 1)
		return err
	case nil:
		atomic.StoreInt32(&t.cnt, 0)
		return nil
	default:
		//不知道什么错误的时候 可以换一个个
		//超时错误 可能偶发的 再试试
		return err
	}
}
