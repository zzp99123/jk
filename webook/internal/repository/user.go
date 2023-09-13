package repository

import (
	"context"
	"database/sql"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository/cache"
	"goFoundation/webook/internal/repository/dao"
	"time"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrNotFound           = dao.ErrNotFound
)

type UserRepositoryIF interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, e string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Update(ctx context.Context, u domain.User) error
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}
type userRepository struct {
	dao   dao.UserDaoIF
	cache cache.UsersCacheIF
}

func NewUserService(dao dao.UserDaoIF, cache cache.UsersCacheIF) UserRepositoryIF {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

// 注册
func (r *userRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

// 登录对比邮箱
func (r *userRepository) FindByEmail(ctx context.Context, e string) (domain.User, error) {
	//return r.dao.FindByEmail(ctx, e)
	u, err := r.dao.FindByEmail(ctx, e)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}

// 查找
func (r *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	//首先先查redis
	//查找成功，缓存命中 get找到了用户 不用在查数据库了
	//查找成功，缓存未命中 get没找到用户  set到了用户
	//未找到用户
	res, err := r.cache.Get(ctx, id)
	if err == nil {
		//有数据
		return res, nil
	}
	//再查数据库 查找成功，缓存未命中
	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	ue := r.entityToDomain(u)
	//go func() {
	//	_ = r.cache.Set(ctx, ue)
	//}()
	_ = r.cache.Set(ctx, ue)
	return ue, nil
}

// 更改
func (r *userRepository) Update(ctx context.Context, u domain.User) error {
	err := r.dao.Update(ctx, r.domainToEntity(u))
	if err != nil {
		return err
	}
	return nil
}

// 查找手机号
func (r *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	ok, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(ok), nil
}

// domain.user转dao
func (r *userRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{
			u.Email,
			u.Email != "",
		},
		Phone: sql.NullString{
			u.Phone,
			u.Phone != "",
		},
		Password: u.Password,
		Birthday: sql.NullInt64{
			u.Birthday.UnixMilli(),
			!u.Birthday.IsZero(),
		},
		Nickname: sql.NullString{
			u.Nickname,
			u.Nickname != "",
		},
		PersonalProfile: sql.NullString{
			u.PersonalProfile,
			u.PersonalProfile != "",
		},
	}
}

// dao.user转domain
func (r *userRepository) entityToDomain(u dao.User) domain.User {
	var birthday time.Time
	if u.Birthday.Valid {
		birthday = time.UnixMilli(u.Birthday.Int64)
	}
	return domain.User{
		Id:              u.Id,
		Email:           u.Email.String,
		Password:        u.Password,
		Phone:           u.Phone.String,
		Nickname:        u.Nickname.String,
		PersonalProfile: u.PersonalProfile.String,
		Birthday:        birthday,
		Ctime:           time.UnixMilli(u.Ctime),
	}
}
