package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"math/rand"
	"time"
)

type lo struct {
	load int64
}

func (r *lo) random() {
	res := cron.New(cron.WithSeconds())
	res.AddFunc("@every 1s", func() {
		r.load = int64(rand.Intn(100))
		fmt.Println(r.load)

	})
	res.Start()
	time.Sleep(time.Second * 10)
	res.Stop()
}
func main() {
	var l lo
	l.random()
}
