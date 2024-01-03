package job

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"goFoundation/webook/pkg/logger"
	"strconv"
	"time"
)

type BuilderJob struct {
	l      logger.LoggerV1
	pro    *prometheus.SummaryVec
	tracer trace.Tracer //接入
}

func NewBuilderJob(l logger.LoggerV1) *BuilderJob {
	pro := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: "geekbang_daming",
		Subsystem: "webook",
		Help:      "统计定时任务的执行情况",
		Name:      "cron_job",
	}, []string{
		"name", "success",
	})
	prometheus.MustRegister(pro)
	return &BuilderJob{
		l:      l,
		pro:    pro,
		tracer: otel.GetTracerProvider().Tracer("webook/internal/job"),
	}
}

// 我们利用一个 Builder 模式来构造各种 cron.Job。
func (b *BuilderJob) Build(job Job) cron.Job {
	name := job.Name()
	return cronJobFuncAdapter(func() error {
		_, span := b.tracer.Start(context.Background(), name)
		defer span.End()
		start := time.Now()
		b.l.Info("任务开始",
			logger.String("job", name))
		var success bool
		defer func() {
			b.l.Info("任务结束",
				logger.String("job", name))
			duration := time.Since(start).Milliseconds()
			b.pro.WithLabelValues(name,
				strconv.FormatBool(success)).Observe(float64(duration))
		}()
		err := job.Run()
		success = err == nil
		if err != nil {
			span.RecordError(err)
			b.l.Error("运行任务失败", logger.Error(err),
				logger.String("job", name))
		}
		return nil
	})
}

type cronJobFuncAdapter func() error

func (c cronJobFuncAdapter) Run() {
	_ = c()
}
