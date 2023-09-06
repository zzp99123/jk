package service

import (
	"context"
	"errors"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type ServiceUser struct {
	r *repository.RepositoryUser
}

var (
	ErruserDuplicateEmail   = repository.ErruserDuplicateEmail
	ErrInvaildUserOrPasswod = errors.New("邮箱/密码不正确")
	ErrNotFound             = repository.ErrNotFound
)

func NewUserService(r *repository.RepositoryUser) *ServiceUser {
	return &ServiceUser{
		r: r,
	}
}

// 注册逻辑
func (s *ServiceUser) SignUp(ctx context.Context, u domain.User) error {
	//密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return s.r.Create(ctx, u)
}

// 登录逻辑
func (s *ServiceUser) Login(ctx context.Context, e string, p string) (domain.User, error) {
	//先找邮箱 在对比密码
	u, err := s.r.FindByEmail(ctx, e)
	//对比邮箱
	if err == ErrNotFound {
		return domain.User{}, ErrInvaildUserOrPasswod
	}
	if err != nil {
		return domain.User{}, err
	}
	//对比密码
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return domain.User{}, ErrInvaildUserOrPasswod
	}
	return u, nil

}

// 查找
func (s *ServiceUser) Profile(ctx context.Context, id int64) (domain.User, error) {
	return s.r.FindById(ctx, id)
}

// 修改
func (s *ServiceUser) Edit(ctx context.Context, user domain.User) error {
	return s.r.Upate(ctx, user)
}
