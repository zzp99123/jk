package tengxunyun

import (
	"context"
	"fmt"
	"github.com/ecodeclub/ekit"
	"github.com/ecodeclub/ekit/slice"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
	mysms "goFoundation/webook/internal/service/sms"
	"goFoundation/webook/pkg/ratelimit"
)

type Service struct {
	appId    *string
	signName *string
	client   *sms.Client
	limiter  ratelimit.Limiter
	//vector   *prometheus.CounterVec //如果你要监控错误码 需要深入到各个细节去弄
}

func NewService(client *sms.Client, appId string,
	signName string, limiter ratelimit.Limiter) *Service {
	//vec := prometheus.NewCounterVec(prometheus.CounterOpts{
	//	Namespace: "geekbang_daming",
	//	Subsystem: "webook",
	//	Name:      "sms_resp_code",
	//	Help:      "统计 SMS 服务错误码",
	//}, []string{"code"})
	//prometheus.MustRegister(vec)
	return &Service{
		client:   client,
		appId:    ekit.ToPtr[string](appId),
		signName: ekit.ToPtr[string](signName),
		limiter:  limiter,
		//vector:   vec,
	}
}

// Send 一个是 []*string
// 一个是 string，json 串
// biz 直接代表的就是 tplId
func (s *Service) Send(ctx context.Context,
	biz string, args []string, numbers ...string) error {
	//这样写 容易屎山 并且代码沉鱼 所以可以用装饰器
	//key := "sms:tencent"
	//l, err := s.limiter.Limit(ctx, key)
	//if err != nil {
	//	//系统错误
	//	//如果第三方很弱的话就限流  如果第三方很强的话就不限流
	//	return fmt.Errorf("短信服务判断是否限流出现问题，%w", err)
	//}
	//if l {
	//	return fmt.Errorf("触发了限流")
	//}
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](biz)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = s.toStringPtrSlice(args)
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		//s.vector.WithLabelValues(*status.Code)
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s ", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) SendV1(ctx context.Context, tplId string, args []mysms.NamedArg, numbers ...string) error {
	req := sms.NewSendSmsRequest()
	req.SmsSdkAppId = s.appId
	req.SignName = s.signName
	req.TemplateId = ekit.ToPtr[string](tplId)
	req.PhoneNumberSet = s.toStringPtrSlice(numbers)
	req.TemplateParamSet = slice.Map[mysms.NamedArg, *string](args, func(idx int, src mysms.NamedArg) *string {
		return &src.Val
	})
	resp, err := s.client.SendSms(req)
	if err != nil {
		return err
	}
	for _, status := range resp.Response.SendStatusSet {
		if status.Code == nil || *(status.Code) != "Ok" {
			return fmt.Errorf("发送短信失败 %s, %s ", *status.Code, *status.Message)
		}
	}
	return nil
}

func (s *Service) toStringPtrSlice(src []string) []*string {
	return slice.Map[string, *string](src, func(idx int, src string) *string {
		return &src
	})
}
