package time

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	//1秒1执行
	res := time.NewTicker(time.Second)
	//避免gorountion泄露
	defer res.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			goto end
		case now := <-res.C:
			t.Log(now)
		}

	}
end:
	t.Log("退出循环")
}
func TestTimer(t *testing.T) {
	//定时器
	//tm := time.NewTimer(time.Second)
	//defer tm.Stop()
	////fatal error: all goroutines are asleep - deadlock! 死锁的意思 用groution
	//go func() {
	//	for now := range tm.C {
	//		t.Log(now.Unix())
	//	}
	//}()
	tm := time.NewTicker(time.Second)
	defer tm.Stop()
	go func() {
		for {
			load := int64(rand.Intn(100))
			fmt.Println(load)
		}
	}()
	time.Sleep(time.Second * 1)
}
