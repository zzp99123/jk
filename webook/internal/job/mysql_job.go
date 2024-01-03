package job

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/pkg/logger"
	"golang.org/x/sync/semaphore"
	"net/http"
	"time"
)

type Executor interface {
	Name() string
	// Exec ctx 是整个任务调度的上下文
	// 当从 ctx.Done 有信号的时候，就需要考虑结束执行
	// 具体实现来控制
	// 真正去执行一个任务
	Exec(ctx context.Context, j domain.CronJob) error
}
type HttpExecutor struct {
}

func (h *HttpExecutor) Name() string {
	return "http"
}

func (h *HttpExecutor) Exec(ctx context.Context, j domain.CronJob) error {
	type Config struct {
		Endpoint string
		Method   string
	}
	var cfg Config
	err := json.Unmarshal([]byte(j.Cfg), &cfg)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(cfg.Method, cfg.Endpoint, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if resp.StatusCode != http.StatusOK {
		return errors.New("执行失败")
	}
	return nil
}

type ExecutorJob struct {
	funcs map[string]func(ctx context.Context, j domain.CronJob) error
}

func NewExecutorJob() *ExecutorJob {
	return &ExecutorJob{
		funcs: make(map[string]func(ctx context.Context, j domain.CronJob) error),
	}
}
func (e *ExecutorJob) Name() string {
	return "ranking"
}
func (e *ExecutorJob) RegisterFunc(name string, fn func(ctx context.Context, j domain.CronJob) error) {
	e.funcs[name] = fn
}
func (e *ExecutorJob) Exec(ctx context.Context, j domain.CronJob) error {
	res, ok := e.funcs[j.Name]
	if !ok {
		return fmt.Errorf("未知任务，你是否注册了？ %s", j.Name)
	}
	return res(ctx, j)
}

type MysqlJob struct {
	svc     service.MysqlJobService
	l       logger.LoggerV1
	ex      map[string]Executor
	limiter *semaphore.Weighted
}

func NewMysqlJob(svc service.MysqlJobService, l logger.LoggerV1) *MysqlJob {
	return &MysqlJob{
		svc: svc,
		l:   l,
		ex:  make(map[string]Executor),
		//通过信号量来限制并行的goroutine数量，达到最大的maxWorkers数量
		limiter: semaphore.NewWeighted(200),
	}
}
func (m *MysqlJob) RegisterExecutor(exec Executor) {
	m.ex[exec.Name()] = exec
}
func (m *MysqlJob) Preempt(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			// 退出调度循环
			return ctx.Err()
		}
		//Acquire将会阻塞，直到其中一个goroutine执行完成，释放出信号量 拿令牌
		err := m.limiter.Acquire(ctx, 1)
		if err != nil {
			return err
		}
		// 一次调度的数据库查询时间
		pctx, cancel := context.WithTimeout(context.Background(), time.Second)
		cancel()
		//抢占
		res, err := m.svc.Preempt(pctx)
		//抢占失败
		if err != nil {
			// 你不能 return
			// 你要继续下一轮
			m.l.Error("抢占任务失败", logger.Error(err))
		}
		//这个执行器是 grpc 还是http
		exec, ok := m.ex[res.Executor] //这给了个标识
		if !ok {
			//debug的时候终端
			//线上就继续
			m.l.Error("没找到对应的执行器", logger.String("executor", res.Executor))
			continue
		}
		go func() {
			//释放权值为n的信号量。执行完毕 释放令牌
			m.limiter.Release(1)
			defer func() {
				//抢占成功后不能一直抢啊 得释放
				err1 := res.CancelFunc()
				if err1 != nil {
					m.l.Error("释放失败", logger.Error(err1), logger.Int64("jid", res.Id))
				}
			}()
			// 异步执行，不要阻塞主调度循环
			// 执行完毕之后
			// 这边要考虑超时控制，任务的超时控制
			err1 := exec.Exec(ctx, res)
			if err1 != nil {
				//也可以在这重试
				m.l.Error("任务执行失败", logger.Error(err1))
			}
			// 你要不要考虑下一次调度？
			ctx1, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			err1 = m.svc.ResetNextTime(ctx1, res)
			if err1 != nil {
				m.l.Error("设置下一次执行时间失败", logger.Error(err1))
			}
		}()
	}

}
