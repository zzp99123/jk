package repository

import (
	"context"
	"database/sql"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/repository/dao"
	"time"
)

type RepositoryUser struct {
	dao *dao.DaoUser
}

var (
	ErruserDuplicateEmail = dao.ErruserDuplicateEmail
	ErrNotFound           = dao.ErrNotFound
)

func NewUserService(dao *dao.DaoUser) *RepositoryUser {
	return &RepositoryUser{
		dao: dao,
	}
}

// 注册
func (r *RepositoryUser) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}

// 登录对比邮箱
func (r *RepositoryUser) FindByEmail(ctx context.Context, e string) (domain.User, error) {
	//return r.dao.FindByEmail(ctx, e)
	u, err := r.dao.FindByEmail(ctx, e)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}

// 查找
func (r *RepositoryUser) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), err
}

// 更改
func (r *RepositoryUser) Upate(ctx context.Context, u domain.User) error {
	err := r.dao.Upate(ctx, r.domainToEntity(u))
	if err != nil {
		return err
	}
	return nil
}

// domain.user转dao
func (r *RepositoryUser) domainToEntity(u domain.User) dao.User {
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
func (r *RepositoryUser) entityToDomain(u dao.User) domain.User {
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
