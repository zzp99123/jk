package ioc

import (
	"context"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/job"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/pkg/logger"
	"time"
)

func InitScheduler(svc service.MysqlJobService, local *job.ExecutorJob, l logger.LoggerV1) *job.MysqlJob {
	res := job.NewMysqlJob(svc, l)
	res.RegisterExecutor(local)
	return res
}
func InitLocalFuncExecutor(svc service.RangingService) *job.ExecutorJob {
	res := job.NewExecutorJob()
	res.RegisterFunc("ranking", func(ctx context.Context, j domain.CronJob) error {
		ctx, cancl := context.WithTimeout(context.Background(), time.Second)
		defer cancl()
		return svc.TopN(ctx)
	})
	return res
}
