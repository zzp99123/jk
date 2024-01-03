// 装饰器 就是我腾讯服务需要那个限流 我阿里也需要限流 不可能每一个都写一边if else 所以我只需要在这个接口的前后写这些我想要的东西
package ratelimit

import (
	"context"
	"fmt"
	"goFoundation/webook/internal/service/sms"
	"goFoundation/webook/pkg/ratelimit"
)

type ServiceRatelimit struct {
	svc     sms.Service
	limiter ratelimit.Limiter
}

func NewServiceRatelimit(svc sms.Service, limiter ratelimit.Limiter) sms.Service {
	return &ServiceRatelimit{
		svc:     svc,
		limiter: limiter,
	}
}
func (s *ServiceRatelimit) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	//装饰的地方
	key := "sms:tencent"
	l, err := s.limiter.Limit(ctx, key)
	if err != nil {
		//系统错误
		//如果第三方很弱的话就限流  如果第三方很强的话就不限流
		return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)
	}
	if l {
		return fmt.Errorf("触发了限流")
	}
	err = s.svc.Send(ctx, tpl, args, numbers...)
	//装饰的地方
	return err
}
