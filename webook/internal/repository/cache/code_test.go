package cache

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"goFoundation/webook/internal/repository/cache/redismocks"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	testCases := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		ctx     context.Context
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "存储成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				c := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(0))
				c.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).Return(res)
				return c
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "发送验证码太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				c := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-1))
				c.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).Return(res)
				return c
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				c := redismocks.NewMockCmdable(ctrl)
				res := redis.NewCmd(context.Background())
				res.SetVal(int64(-10))
				c.EXPECT().Eval(gomock.Any(), luaSetCode, []string{"phone_code:login:152"}, []any{"123456"}).Return(res)
				return c
			},
			ctx:     context.Background(),
			biz:     "login",
			phone:   "152",
			code:    "123456",
			wantErr: errors.New("系统错误"),
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := v.mock(ctrl)
			repo := NewCodeCache(c)
			err := repo.Set(v.ctx, v.biz, v.phone, v.code)
			assert.Equal(t, v.wantErr, err)
		})
	}
}
