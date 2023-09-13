package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository"
	repomocks "goFoundation/webook/internal/repository/mocks"
	"golang.org/x/crypto/bcrypt"
	"testing"
	"time"
)

func Test_userService_Login(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name     string
		mock     func(curl *gomock.Controller) repository.UserRepositoryIF
		email    string
		password string
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "登录成功",
			mock: func(curl *gomock.Controller) repository.UserRepositoryIF {
				r := repomocks.NewMockUserRepositoryIF(curl)
				r.EXPECT().FindByEmail(gomock.Any(), "1234@qq.com").Return(domain.User{
					Email: "1234@qq.com",
					//需要加密
					Password: "$2a$10$cIj/hE5jwJPbVD3Cq30Efu6YnjHdaenLh34fotyYnSh.ONjugyAEa",
					Phone:    "17612422528",
					Ctime:    now,
				}, nil)
				return r
			},
			email:    "1234@qq.com",
			password: "hello#world123",
			wantUser: domain.User{
				Email: "1234@qq.com",
				//需要加密
				Password: "$2a$10$cIj/hE5jwJPbVD3Cq30Efu6YnjHdaenLh34fotyYnSh.ONjugyAEa",
				Phone:    "17612422528",
				Ctime:    now,
			},
			wantErr: nil,
		},
		{
			name: "用户没找到",
			mock: func(curl *gomock.Controller) repository.UserRepositoryIF {
				r := repomocks.NewMockUserRepositoryIF(curl)
				r.EXPECT().FindByEmail(gomock.Any(), "123@qq.com").Return(domain.User{}, repository.ErrNotFound)
				return r
			},
			email:    "123@qq.com",
			password: "hello#world123",
			wantErr:  ErrInvaildUserOrPassword,
		},
		{
			name: "密码错误",
			mock: func(curl *gomock.Controller) repository.UserRepositoryIF {
				r := repomocks.NewMockUserRepositoryIF(curl)
				r.EXPECT().FindByEmail(gomock.Any(), "1234@qq.com").Return(domain.User{
					Email: "1234@qq.com",
					//需要加密
					Password: "$2a$10$cIj/hE5jwJPbVD3Cq30Efu6YnjHdaenLh34fotyYnSh.ONjugyAEa",
					Phone:    "17612422528",
					Ctime:    now,
				}, nil)
				return r
			},
			email:    "1234@qq.com",
			password: "hello#world1",

			wantErr: ErrInvaildUserOrPassword,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			curl := gomock.NewController(t)
			defer curl.Finish()
			svc := NewUserService(tc.mock(curl))
			u, err := svc.Login(context.Background(), tc.email, tc.password)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, u)

		})
	}
}

// 密码加密
func TestEncrypted(t *testing.T) {
	res, err := bcrypt.GenerateFromPassword([]byte("hello#world123"), bcrypt.DefaultCost)
	if err == nil {
		t.Log(string(res))
	}
}
