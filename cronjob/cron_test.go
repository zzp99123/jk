package cronjob

import (
	"fmt"
	cron "github.com/robfig/cron/v3"
	"testing"
	"time"
)

func TestCronJob(t *testing.T) {
	c := cron.New(cron.WithSeconds())
	//c.AddJob("@every 1s", Myjob{})
	c.AddFunc("@every 3s", func() {
		t.Log("func执行了")
		time.Sleep(time.Second * 12)
		t.Log("func10秒后执行了")
	})
	// s 就是秒, m 就是分钟, h 就是小时，d 就是天
	c.Start()
	//模仿执行
	time.Sleep(time.Second * 10)
	//只是通知你结束
	t.Log("通知结束")
	// 发出停止信号，expr 不会调度新的任务，但是也不会中断已经调度了的任务
	s := c.Stop()
	// 这一句会阻塞，等到所有已经调度（正在运行的）结束，才会返回
	<-s.Done() //这才是真正的结束
	t.Log("彻底结束")
}

type Myjob struct {
}

func (m Myjob) Run() {
	fmt.Println("add运行了")
}
