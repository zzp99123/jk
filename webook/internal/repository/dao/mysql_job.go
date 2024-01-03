// 分布式任务调度系统
package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type MysqlJobDao interface {
	Preempt(ctx context.Context) (CronJob, error)
	UpdateUtime(ctx context.Context, id int64) error
	Release(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	ResetNextTime(ctx context.Context, id int64, next time.Time) error
}
type mysqlJobDao struct {
	db *gorm.DB
}

func NewMysqlJobDao(db *gorm.DB) MysqlJobDao {
	return &mysqlJobDao{
		db: db,
	}
}

// 下一次执行时间设置
func (m *mysqlJobDao) ResetNextTime(ctx context.Context, id int64, next time.Time) error {
	return m.db.WithContext(ctx).Model(&CronJob{}).Where("id", id).Updates(map[string]any{
		"next_time": next.UnixMilli(),
	}).Error
}

// 续约
func (m *mysqlJobDao) UpdateUtime(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Model(&CronJob{}).Where("id", id).Updates(map[string]any{
		"utime": time.Now().UnixMilli(),
	}).Error
}

// 抢占
func (m *mysqlJobDao) Preempt(ctx context.Context) (CronJob, error) {
	for {
		now := time.Now()
		var c CronJob
		// 1. 一次拉一批，我一次性取出 100 条来，然后，我随机从某一条开始，向后开始抢占
		// 2. 我搞个随机偏移量，0-100 生成一个随机偏移量。兜底：第一轮没查到，偏移量回归到 0
		// 3. 我搞一个 id 取余分配，status = ? AND next_time <=? AND id%10 = ? 兜底：不加余数条件，取next_time 最老的
		//我先找什么样的任务可以抢
		err := m.db.WithContext(ctx).Where("next_time = ? And Status = ? ", now, jobStatusWaiting).First(&c).Error
		if err != nil {
			return CronJob{}, err
		}
		// 你找到了，可以被抢占的
		// 找到之后你要干嘛？你要抢占
		res := m.db.WithContext(ctx).Where("Version = ? AND id = ?", c.Version, c.Id).Model(&c).
			Updates(map[string]any{
				"status":  jobStatusRunning,
				"utime":   now,
				"version": c.Version + 1,
			})
		if res.Error != nil {
			return CronJob{}, err
		}
		if res.RowsAffected == 0 {
			continue
		}
		return c, nil
	}
}

// 释放
func (m *mysqlJobDao) Release(ctx context.Context, id int64) error {
	// 这里有一个问题。你要不要检测 status 或者 version?
	// WHERE version = ?
	// 要。你们的作业记得修改
	return m.db.WithContext(ctx).Where("id", id).
		Updates(map[string]any{
			"status": jobStatusPaused,
			"utime":  time.Now().UnixMilli(),
		}).Error

}

// 暂停
func (m *mysqlJobDao) Stop(ctx context.Context, id int64) error {
	return m.db.WithContext(ctx).Where("id", id).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  time.Now().UnixMilli(),
	}).Error
}

type CronJob struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`
	// 另外一个问题，定时任务，我怎么知道，已经到时间了呢？
	// next_time <= now 这样一个查询条件
	NextTime int64 `gorm:"index"` // 下一次被调度的时间
	Cron     string
	Executor string
	Name     string `gorm:"unique"`
	Cfg      string
	Status   int
	Version  int
	Ctime    int64
	UTime    int64
}

const (
	jobStatusWaiting = iota //暂停
	jobStatusRunning        //抢占
	jobStatusPaused         //释放
	// 还可以有别的取值
)
