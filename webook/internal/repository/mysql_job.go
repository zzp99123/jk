package repository

import (
	"context"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository/dao"
	"time"
)

type MysqlJobRepository interface {
	Preempt(ctx context.Context) (domain.CronJob, error)
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	ResetNextTime(ctx context.Context, id int64, next time.Time) error
}
type mysqlJobRepository struct {
	dao dao.MysqlJobDao
}

func NewMysqlJobRepository(dao dao.MysqlJobDao) MysqlJobRepository {
	return &mysqlJobRepository{
		dao: dao,
	}
}

// 续约
func (m *mysqlJobRepository) ResetNextTime(ctx context.Context, id int64, next time.Time) error {
	return m.dao.ResetNextTime(ctx, id, next)
}

// 续约
func (m *mysqlJobRepository) UpdateUtime(ctx context.Context, id int64) error {
	return m.dao.UpdateUtime(ctx, id)
}

// 释放
func (m *mysqlJobRepository) Release(ctx context.Context, id int64) error {
	return m.dao.Release(ctx, id)
}

// 暂停
func (m *mysqlJobRepository) Stop(ctx context.Context, id int64) error {
	return m.dao.Stop(ctx, id)
}

// 抢占
func (m *mysqlJobRepository) Preempt(ctx context.Context) (domain.CronJob, error) {
	res, err := m.dao.Preempt(ctx)
	if err != nil {
		return domain.CronJob{}, err
	}
	return domain.CronJob{
		Cfg:      res.Cfg,
		Name:     res.Name,
		Id:       res.Id,
		Executor: res.Executor,
	}, nil

}
