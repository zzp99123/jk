package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

type PrometheusMetric struct {
	Namespace  string
	Subsystem  string
	Name       string
	Help       string
	InstanceId string
}

func (p *PrometheusMetric) Build() gin.HandlerFunc {
	// pattern 是指你命中的路由
	// 是指你的 HTTP 的 status
	// path /detail/1
	label := []string{"method", "pattern", "status"} //数量不能太大 否则Prometheus堵塞
	summary := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		//前三个不能用- 可以用下划线_
		Namespace: p.Namespace,
		Subsystem: p.Subsystem,
		Name:      p.Name + "_resp_time",
		Help:      p.Help,
		ConstLabels: map[string]string{
			"instance_id": p.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	}, label)
	prometheus.MustRegister(summary)
	//当前请求活跃数量gauge
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: p.Namespace,
		Subsystem: p.Subsystem,
		Name:      p.Name + "_active_req",
		Help:      p.Help,
		ConstLabels: map[string]string{
			"instance_id": p.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		start := time.Now()
		gauge.Inc() //加1
		defer func() {
			duration := time.Since(start)
			gauge.Dec() //减1
			//如果404
			pattern := ctx.FullPath() //获取路由
			if pattern == "" {
				pattern = "unknown"
			}
			method := ctx.Request.Method
			status := ctx.Writer.Status()
			summary.WithLabelValues(method, pattern, strconv.Itoa(status)).Observe(float64(duration.Milliseconds()))
		}()
		// 你最终就会执行到业务里面
		ctx.Next()
	}
}
