package redisx

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"net"
	"strconv"
	"time"
)

type PrometheusRedisx struct {
	//client redis.Client
	vec *prometheus.SummaryVec
}

func NewPrometheusRedisx(top prometheus.SummaryOpts) *PrometheusRedisx {
	// key_exist 是否命中缓存
	vec := prometheus.NewSummaryVec(top, []string{"cmd", "key_exist"})
	prometheus.MustRegister(vec)
	return &PrometheusRedisx{
		vec: vec,
	}
}

//	func Use(client *redis.Client) {
//		client.AddHook()
//	}
//
// 跟redis建立链接
func (p *PrometheusRedisx) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// 发送普通命令
func (p *PrometheusRedisx) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {
		startime := time.Now()
		var err error
		//cmd.name方法是用来看我用的是0还是1 0是get方法 1是set方法
		defer func() {
			//biz := ctx.Value("user") //获取深入监控方法的
			keyExist := err == redis.Nil //判断缓存命中
			p.vec.WithLabelValues(cmd.Name(), strconv.FormatBool(keyExist)).Observe(float64(time.Since(startime).Milliseconds()))
		}()
		err = next(ctx, cmd)
		return err
	}
}

//使用Pipeline功能

func (p *PrometheusRedisx) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}
