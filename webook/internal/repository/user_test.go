package repository

import (
	"context"
	"database/sql"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository/cache"
	cachemocks "goFoundation/webook/internal/repository/cache/mocks"
	"goFoundation/webook/internal/repository/dao"
	daomocks "goFoundation/webook/internal/repository/dao/mocks"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	now := time.Now().UnixMilli()
	nows := time.UnixMilli(now)
	testCases := []struct {
		name     string
		mock     func(ctrl *gomock.Controller) (dao.UserDao, cache.UsersCache)
		ctx      context.Context
		id       int64
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UsersCache) {
				d := daomocks.NewMockUserDaoIF(ctrl)
				c := cachemocks.NewMockUsersCacheIF(ctrl)
				// int64(123) int类型需要变换int64类型
				//模拟缓存未命中
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				c.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Phone:    "15212345678",
					Ctime:    nows}).Return(nil)

				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{
					Id: 123,
					Email: sql.NullString{
						String: "123@qq.com",
						Valid:  true,
					},
					Password: "123456",
					Phone: sql.NullString{
						String: "15212345678",
						Valid:  true,
					},
					Ctime: now,
					UTime: now,
				}, nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "15212345678",
				Ctime:    nows,
			},
		},
		{
			name: "查找成功，缓存命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UsersCache) {
				d := daomocks.NewMockUserDaoIF(ctrl)
				c := cachemocks.NewMockUsersCacheIF(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{
					Id:       123,
					Email:    "123@qq.com",
					Password: "123456",
					Phone:    "15212345678",
					Ctime:    nows,
				}, nil)
				return d, c
			},
			ctx: context.Background(),
			id:  123,
			wantUser: domain.User{
				Id:       123,
				Email:    "123@qq.com",
				Password: "123456",
				Phone:    "15212345678",
				Ctime:    nows,
			},
		},
		{
			name: "没找到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDao, cache.UsersCache) {
				d := daomocks.NewMockUserDaoIF(ctrl)
				c := cachemocks.NewMockUsersCacheIF(ctrl)
				c.EXPECT().Get(gomock.Any(), int64(123)).Return(domain.User{}, cache.ErrKeyNotExist)
				d.EXPECT().FindById(gomock.Any(), int64(123)).Return(dao.User{}, ErrNotFound)
				return d, c
			},
			ctx:     context.Background(),
			id:      123,
			wantErr: ErrNotFound,
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			u, c := v.mock(ctrl)
			repo := NewUserService(u, c)
			d, err := repo.FindById(v.ctx, v.id)
			assert.Equal(t, v.wantErr, err)
			assert.Equal(t, v.wantUser, d)
			//如果是并发用time.sleep休眠一下
		})
	}
}
