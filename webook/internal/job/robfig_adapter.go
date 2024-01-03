package job

import (
	"github.com/prometheus/client_golang/prometheus"
	"goFoundation/webook/pkg/logger"
	"time"
)

type RobfigAdapterJob struct {
	job Job
	l   logger.LoggerV1
	pro prometheus.Summary
}

func NewRobfigAdapterJob(job Job, l logger.LoggerV1) *RobfigAdapterJob {
	pro := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "cron_job",
		ConstLabels: map[string]string{
			"name": job.Name(),
		},
	})
	prometheus.MustRegister(pro)
	return &RobfigAdapterJob{
		job: job,
		l:   l,
		pro: pro,
	}
}
func (r *RobfigAdapterJob) Run() {
	startTime := time.Now()
	defer func() {
		r.pro.Observe(float64(time.Since(startTime).Milliseconds()))
	}()
	err := r.job.Run()
	if err != nil {
		r.l.Error("运行失败", logger.Error(err), logger.String("job", r.job.Name()))
	}
}
