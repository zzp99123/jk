package ioc

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	"goFoundation/webook/internal/repository/cache"
	"goFoundation/webook/pkg/redisx"
)

// 这是为了更深入的监控到各个方法 如果需要也是得每个模块都得定义一遍
func InitUserCache(client *redis.Client) cache.UsersCache {
	client.AddHook(redisx.NewPrometheusRedisx(prometheus.SummaryOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Name:      "gin_http",
		Help:      "统计 GIN 的 HTTP 接口",
		ConstLabels: map[string]string{
			"biz": "user",
		},
	}))
	panic("你别调用")
}
