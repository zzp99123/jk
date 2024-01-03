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
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrNotFound           = gorm.ErrRecordNotFound
)

// ErrDataNotFound 通用的数据没找到
var ErrDataNotFound = gorm.ErrRecordNotFound

type UserDao interface {
	Insert(ctx context.Context, u User) error
	FindByEmail(ctx context.Context, e string) (User, error)
	FindByPhone(ctx context.Context, phone string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Update(ctx context.Context, u User) error
	FindByWechat(ctx context.Context, openId string) (User, error)
}

type userDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &userDao{
		db: db,
	}
}

// 注册
func (d *userDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.UTime = now
	err := d.db.WithContext(ctx).Create(&u).Error
	//检验重复邮箱注册
	if me, ok := err.(*mysql.MySQLError); ok {
		const uniqueIndexErrorNo uint16 = 1062
		if me.Number == uniqueIndexErrorNo {
			return ErrUserDuplicateEmail
		}
	}
	return err
}

// 登录并且先找邮箱
func (d *userDao) FindByEmail(ctx context.Context, e string) (User, error) {
	var u User
	//找到数据库存储的邮箱并返回给上一层
	err := d.db.WithContext(ctx).Where("email = ?", e).First(&u).Error
	return u, err
}

// 查找手机号
func (d *userDao) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("phone = ?", phone).First(&u).Error
	return u, err
}

// 微信扫码
func (d *userDao) FindByWechat(ctx context.Context, openId string) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("wechat_open_id = ?", openId).First(&u).Error
	return u, err
}

// 查找id
func (d *userDao) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := d.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	return u, err
}

// 更改
func (d *userDao) Update(ctx context.Context, u User) error {
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
	UTime int64
	//扫码登录的校验用的值
	WechatOpenId  sql.NullString `gorm:"unique"`
	WechatUnionId sql.NullString `gorm:"unique"`

	Birthday        sql.NullInt64  //生日
	PersonalProfile sql.NullString `gorm:"type=varchar(1024)"` //个人简介
	Nickname        sql.NullString //昵称
}
