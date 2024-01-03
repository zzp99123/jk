// Context 接口核心 API 有四个：
package context

import (
	"context"
	"testing"
	"time"
)

//type Key1 struct{}

// context.Value：取值。非常常用
func TestContext(t *testing.T) {
	ctx1 := context.WithValue(context.Background(), "key1", "val1")
	val := ctx1.Value("key1") //官方建议用结构题
	t.Log(val)
	ctx2 := context.WithValue(context.Background(), "key2", "val2")
	val = ctx1.Value("key1")
	t.Log(val)
	val = ctx2.Value("key2")
	t.Log(val)
	ctx1 = context.WithValue(context.Background(), "key1", "val1-1")
	val = ctx1.Value("key1")
	t.Log(val)
}

// // Done：返回一个 channel，一般用于监听Context 实例的信号。比如说过期，或者正常关闭。常用
func TestContext_Cancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	// 为什么一定要 cancel 呢？
	// 防止 goroutine 泄露
	defer cancel()
	// 防止有些人使用了 Done，在等待ctx结束信号
	go func() {
		ch := ctx.Done()
		<-ch
	}()
	ctx1 := context.WithValue(ctx, "key1", "val1-1")
	val := ctx1.Value("key1")
	t.Log(val)
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	cancel()
	// Deadline ：返回过期时间，如果 ok 为 false，说明没有设置过期时间。不常用
	ctx, cancel = context.WithDeadline(ctx, time.Now().Add(time.Second))
	cancel()
}

// Err：返回一个错误用于表达 Context 发生了什么。
func TestContextErr(t *testing.T) {
	//ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	//defer cancel()
	//if ctx.Err() == context.Canceled 	//被取消 {
	//
	//
	//} else if ctx.Err() == context.DeadlineExceeded //超时 {
	//
	//}
}

// 特点：context 的实例之间存在父子关系：
// 当父亲取消或者超时，所有派生的子 context 都被取消或者超时。控制是从上至下的。
func TestContextSub(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	subctx, _ := context.WithCancel(ctx)
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	go func() {
		// 监听 subCtx 结束的信号
		t.Log("等待信号...")
		<-subctx.Done()
		t.Log("收到信号...")
	}()
	time.Sleep(time.Second * 10)
}

// 当找 key 的时候，子 context 先看自己有没有，没有则去祖先里面找。查找则是从下至上的。
func TestContextSubCancel(t *testing.T) {
	ctx, _ := context.WithCancel(context.Background())
	_, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(time.Second)
		cancel()
	}()
	go func() {
		// 监听 subCtx 结束的信号
		t.Log("等待信号...")
		<-ctx.Done()
		t.Log("收到信号...")
	}()
	time.Sleep(time.Second * 10)
}

//func MockIO() {
//	select {
//
//	case <-ctx.Done():
//		// 监听超时 或者用户主动取消
//
//	case <-biz.Signal():
//		// 监听你的正常业务
//	}
//}
