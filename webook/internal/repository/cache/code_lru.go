// 作业实现本地缓存
// 定义一个 CodeCache 接口，将现在的 CodeCache 改名为 CodeRedisCache。
// 提供一个基于本地缓存的 CodeCache 实现。你可以自主决定用什么本地缓存，在这个过程注意体会技术选型要考虑的点。
// 保证单机并发安全，也就是你可以假定这个实现只用在开发环境，或者单机环境下。
package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"sync"
	"time"
)

// 技术选型考虑的点
//  1. 功能性：功能是否能够完全覆盖你的需求。
//  2. 社区和支持度：社区是否活跃，文档是否齐全，
//     以及百度（搜索引擎）能不能搜索到你需要的各种信息，有没有帮你踩过坑
//  3. 非功能性：易用性（用户友好度，学习曲线要平滑），
//     扩展性（如果开源软件的某些功能需要定制，框架是否支持定制，以及定制的难度高不高）
//     性能（追求性能的公司，往往有能力自研）
type CodeLruCache struct {
	l          *lru.Cache    //lru本地缓存
	lock       sync.Mutex    //普通锁，或者说写锁
	rwLock     *sync.RWMutex // 可以多个人加读锁 读写锁
	expiration time.Duration // 我选用的本地缓存，很不幸的是，没有获得过期时间的接口，所以都是自己维持了一个过期时间字段
	maps       sync.Map      //在并发环境中使用的map
}

func NewCodeLruCache(l *lru.Cache, expiration time.Duration) *CodeLruCache {
	return &CodeLruCache{
		l:          l,
		expiration: expiration,
	}
}
func (c *CodeLruCache) Set(ctx context.Context, biz, phone, code string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	//先获取key
	keys := c.key(biz, phone)
	//然后通过get看有没有发送过 如果没有就发送 有先比较时间
	now := time.Now()
	itm, ok := c.l.Get(keys)
	//就说明第一次发
	if !ok {
		//我就添加上
		c.l.Add(keys, codeTime{
			code:  code,
			cnt:   3,
			expir: now.Add(c.expiration),
		})
		return nil
	}
	val, ok := itm.(codeTime)
	if !ok {
		return errors.New("系统错误")
	}
	if val.expir.Sub(now) > time.Minute*9 {
		// 不到一分钟
		return ErrCodeSendTooMany
	}
	//重发
	c.l.Add(keys, codeTime{
		code:  code,
		cnt:   3,
		expir: now.Add(c.expiration),
	})
	return nil
}
func (c *CodeLruCache) Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	keys := c.key(biz, phone)
	val, ok := c.l.Get(keys)
	if !ok {
		//一次都没发过
		return false, ErrKeyNotExist
	}
	itm, ok := val.(codeTime)
	if !ok {
		return false, errors.New("系统错误")
	}
	if itm.cnt <= 0 {
		return false, ErrCodeVerifyTooManyTimes
	}
	itm.cnt--
	return itm.code == expectedCode, nil
}
func (c *CodeLruCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeTime struct {
	cnt   int       // 可验证次数
	code  string    //验证码
	expir time.Time //过期时间

}
