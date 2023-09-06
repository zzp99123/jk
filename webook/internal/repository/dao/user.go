package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErruserDuplicateEmail = errors.New("邮箱冲突")
	ErrNotFound           = gorm.ErrRecordNotFound
)

type DaoUser struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *DaoUser {
	return &DaoUser{
		db: db,
	}
}

// 注册
func (d *DaoUser) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := d.db.WithContext(ctx).Create(&u).Error
	//检验重复邮箱注册
	if me, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrorNo uint16 = 1062
		if me.Number == uniqueIndexErrorNo {
			return ErruserDuplicateEmail
		}
	}
	return err
}

// 登录并且先找邮箱
func (d *DaoUser) FindByEmail(ctx context.Context, e string) (User, error) {
	var u User
	//找到数据库存储的邮箱并返回给上一层
	err := d.db.WithContext(ctx).Where("email = ?", e).First(&u).Error
	return u, err
}

// 查找
func (d *DaoUser) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

// 更改
func (d *DaoUser) Upate(ctx context.Context, u User) error {
	//查找表里的id 去更改其他的数据
	return d.db.Updates(&u).Error
}

type User struct {
	Id       int64          `gorm:"primaryKey,autoIncrement"`
	Email    sql.NullString `gorm:"unique"`
	Password string
	Phone    sql.NullString `gorm:"unique"`
	//创建时间
	Ctime int64
	//更新时间
	Utime           int64
	Birthday        sql.NullInt64  //生日
	PersonalProfile sql.NullString `gorm:"type=varchar(1024)"` //个人简介
	Nickname        sql.NullString //昵称
}
