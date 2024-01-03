package metrics

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/service/oauth2/wechat"
	"time"
)

type PrometheusWechat struct {
	vec prometheus.Summary
	svc wechat.Service
}

func NewprometheusWechat(svc wechat.Service, namespace string,
	subsystem string,
	instanceId string,
	name string) *PrometheusWechat {
	vec := prometheus.NewSummary(prometheus.SummaryOpts{
		Name:      name,
		Namespace: namespace,
		Subsystem: subsystem,
		ConstLabels: map[string]string{
			"instance_id": instanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.95:  0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	})
	prometheus.MustRegister(vec)
	return &PrometheusWechat{
		svc: svc,
		vec: vec,
	}
}

// VerifyCode 因为 AuthURL 过于简单，没有监控的必要
// 定义装饰器把监控加进去
func (p *PrometheusWechat) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	startime := time.Now()
	defer func() {
		p.vec.Observe(float64(time.Since(startime).Milliseconds()))
	}()
	return p.svc.VerifyCode(ctx, code)
}
