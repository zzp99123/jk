package service

import (
	"context"
	"errors"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvaildUserOrPassword = errors.New("邮箱/密码不正确")
)

type UserServiceIF interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, e string, p string) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
	Profile(ctx context.Context, id int64) (domain.User, error)
	Edit(ctx context.Context, user domain.User) error
}
type userService struct {
	r repository.UserRepositoryIF
}

func NewUserService(r repository.UserRepositoryIF) UserServiceIF {
	return &userService{
		r: r,
	}
}

// 邮箱注册逻辑
func (s *userService) SignUp(ctx context.Context, u domain.User) error {
	//密码加密
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return s.r.Create(ctx, u)
}

// 邮箱登录逻辑
func (s *userService) Login(ctx context.Context, e string, p string) (domain.User, error) {
	//先找邮箱 在对比密码
	u, err := s.r.FindByEmail(ctx, e)
	//对比邮箱 用户没找到
	if err == repository.ErrNotFound {
		return domain.User{}, ErrInvaildUserOrPassword
	}
	//对比密码 密码错误
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(p))
	if err != nil {
		return domain.User{}, ErrInvaildUserOrPassword
	}
	return u, nil

}

// 手机号注册登录
func (s *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	res, err := s.r.FindByPhone(ctx, phone)
	if err != repository.ErrNotFound {
		return res, err
	}
	//当到这一步的时候就说明没注册过 这时候就需要注册
	res = domain.User{
		Phone: phone,
	}
	err = s.r.Create(ctx, res)
	//主从延迟
	return s.r.FindByPhone(ctx, phone)
}

// 查找
func (s *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return s.r.FindById(ctx, id)
}

// 修改
func (s *userService) Edit(ctx context.Context, user domain.User) error {
	return s.r.Update(ctx, user)
}
