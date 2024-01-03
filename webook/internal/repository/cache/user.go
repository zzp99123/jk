// 提升性能用户会先从 Redis 里面查询，而后在缓存未命中的 情况下，就会直接从数据库中查询。
package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"goFoundation/webook/internal/domain"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UsersCache interface {
	Get(ctx context.Context, id int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
	Id(id int64) string
}
type usersCache struct {
	client redis.Cmdable
	//过期时间
	expiration time.Duration
}

func NewUsersCache(client redis.Cmdable) UsersCache {
	return &usersCache{
		client,
		time.Minute * 15,
	}
}
func (c *usersCache) Get(ctx context.Context, id int64) (domain.User, error) {
	//这种方法也是深入的监控各个方法 如果要用的话每个方法都需要写一遍
	//ctx = context.WithValue(ctx, "biz", "user")
	//这个是子任务的子任务的监控方法
	//ctx = context.WithValue(ctx, "pattern", "user:info:%d")
	key := c.Id(id)
	res, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(res, &u)
	return u, err

}
func (c *usersCache) Set(ctx context.Context, u domain.User) error {
	//ctx = context.WithValue(ctx, "biz", "user")
	//ctx = context.WithValue(ctx, "pattern", "user:info:%d")
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := c.Id(u.Id)
	return c.client.Set(ctx, key, val, c.expiration).Err()
}
func (c *usersCache) Id(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
