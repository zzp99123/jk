// 添加监控，确保能够掌握 第三方调用的性能数据。
package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"goFoundation/webook/internal/service/sms"
	"time"
)

type PrometheusDecorator struct {
	vec *prometheus.SummaryVec
	svc sms.Service
}

func NewPrometheusDecorator(svc sms.Service) *PrometheusDecorator {
	vec := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "sms_resp_time",
		Help:      "统计 SMS 服务的性能数据",
	}, []string{"biz"})
	prometheus.MustRegister(vec)
	return &PrometheusDecorator{
		vec: vec,
		svc: svc,
	}

}
func (p *PrometheusDecorator) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	startime := time.Now()
	defer func() {
		p.vec.WithLabelValues(biz).Observe(float64(time.Since(startime).Milliseconds()))
	}()
	err := p.svc.Send(ctx, biz, args, numbers...)
	return err
}
