package jwt

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

var (
	AccessKey  = []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0")
	RefreshKey = []byte("95osj3fUu8fo0mlYdDbncXz4VD2igvf0")
)

type RedisJWTHandler struct {
	client       redis.Cmdable
	rtExpiration time.Duration
}

func NewRedisJWTHandler(client redis.Cmdable) Handler {
	return &RedisJWTHandler{
		client:       client,
		rtExpiration: time.Hour * 24 * 7,
	}
}
func (r *RedisJWTHandler) SetLoginToken(ctx *gin.Context, uid int64) error {
	ssid := uuid.New().String()
	err := r.SetJWTToken(ctx, uid, ssid)
	if err != nil {
		return err
	}
	err = r.refJWTToken(ctx, uid, ssid)
	return err
}
func (r *RedisJWTHandler) SetJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
		Ssid:      ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(AccessKey)
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}
func (r *RedisJWTHandler) refJWTToken(ctx *gin.Context, uid int64, ssid string) error {
	claims := RefClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
		Uid:  uid,
		Ssid: ssid,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString(RefreshKey)
	if err != nil {
		return err
	}
	ctx.Header("x-ref-token", tokenStr)
	return nil
}
func (r *RedisJWTHandler) ExtractToken(ctx *gin.Context) string {
	tokenHeader := ctx.GetHeader("Authorization")
	//如果登录过的花 先切割
	segs := strings.Split(tokenHeader, " ")
	if len(segs) > 2 {
		//格式不对登录
		return ""
	}
	return segs[1]
}
func (r *RedisJWTHandler) ClearToken(ctx *gin.Context) error {
	ctx.Header("x-jwt-token", "")
	ctx.Header("x-ref-token", "")
	claims := ctx.MustGet("claims").(*UserClaims)
	return r.client.Set(ctx, fmt.Sprintf("users:ssid:%s", claims.Ssid), "", time.Hour*24*7).Err()
}
func (r *RedisJWTHandler) CheckSession(ctx *gin.Context, ssid string) error {
	res, err := r.client.Exists(ctx, fmt.Sprintf("users:ssid:%s", ssid)).Result()
	//return err
	switch err {
	case redis.Nil:
		return nil
	case nil:
		if res == 0 {
			return err
		}
		return errors.New("session已失效")
	default:
		return err
	}
}

//布隆过滤器主要作用就是查询一个数据，在不在这个二进制的集合中，查询过程如下：
//通过K个哈希函数计算该数据，对应计算出的K个hash值
//通过hash值找到对应的二进制的数组下标
//判断：如果存在一处位置的二进制数据是0，那么该数据不存在。如果都是1，该数据存在集合中。（这里有缺点，下面会讲）
