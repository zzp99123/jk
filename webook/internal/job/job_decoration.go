package job

import (
	"context"
	"math/rand"
	"time"
)

type DecorationJob struct {
	j    MysqlJob
	load int64
}

func (d *DecorationJob) Preempt(ctx context.Context) error {
	tm := time.NewTicker(time.Second)
	defer tm.Stop()
	go func() {
		d.load = int64(rand.Intn(100))
	}()
	err := d.j.Preempt(ctx)
	return err
}
