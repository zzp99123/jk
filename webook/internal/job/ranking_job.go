package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/pkg/logger"
	"sync"
	"time"
)

type RankingJob struct {
	svc       service.RangingService
	timeout   time.Duration
	client    *rlock.Client
	key       string
	l         logger.LoggerV1
	lock      *rlock.Lock
	localLock *sync.Mutex
}

func NewRankingJob(svc service.RangingService, client *rlock.Client, l logger.LoggerV1, timeout time.Duration) *RankingJob {
	return &RankingJob{
		svc:       svc,
		timeout:   timeout,
		l:         l,
		client:    client,
		key:       "rlock:cron_job:ranking",
		localLock: &sync.Mutex{},
	}
}

func (r *RankingJob) Name() string {
	return "ranking"
}

// 如果我们部署了多个实例，那么我们希望，一直都只有一个节点在运行这个榜单任务
// 分布式锁的效果是可以确保整个分布式环境下，只有一个 goroutine 能够拿到锁。
// 按时间调度的，三分钟一次
func (r *RankingJob) Run() error {
	r.localLock.Lock()
	defer r.localLock.Unlock()
	if r.lock == nil {
		// 说明你没拿到锁，你得试着拿锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		// 我可以设置一个比较短的过期时间
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      0,
		}, time.Second)
		if err != nil {
			// 这边没拿到锁，极大概率是别人持有了锁
			return nil
		}
		r.lock = lock
		// 我怎么保证我这里，一直拿着这个锁？？？
		go func() {
			// 自动续约机制
			err1 := lock.AutoRefresh(r.timeout/2, time.Second)
			// 这里说明退出了续约机制
			// 续约失败了怎么办？
			if err1 != nil {
				// 不怎么办
				// 争取下一次，继续抢锁
				r.l.Error("续约失败", logger.Error(err))
			}
			r.localLock.Lock()
			r.lock = nil
			r.localLock.Unlock()
			// lock.Unlock(ctx)
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	return r.svc.TopN(ctx)
}

// 实际上不调用也可以的。因为你关机之后，分布式锁没有人续约，过一会就会有别的节点能够拿到别的分布式锁，继续执行
func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.lock = nil
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}
