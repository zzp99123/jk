package service

import (
	"context"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository"
	"goFoundation/webook/pkg/logger"
	"time"
)

type MysqlJobService interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	UpdateUtime(id int64) error
	//Release(ctx context.Context, id int64) error
	ResetNextTime(ctx context.Context, j domain.CronJob) error
}
type mysqlJobService struct {
	repo            repository.MysqlJobRepository
	refreshInterval time.Duration
	l               logger.LoggerV1
}

func NewMysqlJobService(repo repository.MysqlJobRepository, refreshInterval time.Duration, l logger.LoggerV1) MysqlJobService {
	return &mysqlJobService{
		repo:            repo,
		refreshInterval: refreshInterval,
		l:               l,
	}
}
func (m *mysqlJobService) Preempt(ctx context.Context) (domain.CronJob, error) {
	//抢占
	res, _ := m.repo.Preempt(ctx)
	//抢占失败的话不能直接return出去 这里要用续约机制
	ticker := time.NewTicker(m.refreshInterval)
	go func() {
		for range ticker.C {
			m.UpdateUtime(res.Id)
		}
	}()
	res.CancelFunc = func() error {
		ticker.Stop()
		ctx, cancal := context.WithTimeout(context.Background(), time.Second)
		defer cancal()
		return m.repo.Release(ctx, res.Id)
	}
	//ticker.Stop()
	//err = m.Release(ctx, res.Id)
	//if err != nil {
	//	panic(err)
	//}
	return res, nil
}

// 续约
func (m *mysqlJobService) UpdateUtime(id int64) error {
	ctx, cancal := context.WithTimeout(context.Background(), time.Second)
	defer cancal()
	err := m.repo.UpdateUtime(ctx, id)
	if err != nil {
		//记录日志
		m.l.Error("续约失败", logger.Error(err), logger.Int64("jid", id))
	}
	return nil
}

// 你抢占之后，你一直抢占着吗？你要考虑一个释放的问题
//func (m *mysqlJobService) Release(ctx context.Context, id int64) error {
//	return m.repo.Release(ctx, id)
//}

// 下一次调度 你没有的话就暂停 有的话继续进行下一次
func (m *mysqlJobService) ResetNextTime(ctx context.Context, j domain.CronJob) error {
	next := j.NextTime()
	if next.IsZero() {
		// 没有下一次
		return m.repo.Stop(ctx, j.Id)
	}
	return m.repo.ResetNextTime(ctx, j.Id, next)
}
