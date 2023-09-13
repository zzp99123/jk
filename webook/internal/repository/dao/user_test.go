package dao

import (
	"context"
	"database/sql"
	"errors"
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGORMUserDAO_Insert(t *testing.T) {
	testCases := []struct {
		name    string
		sqlmock func(t *testing.T) *sql.DB
		ctx     context.Context
		u       User
		wantErr error
	}{
		{
			name: "插入成功",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mockRes := sqlmock.NewResult(1, 1)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnResult(mockRes)
				return db
			},
			ctx: context.Background(),
			u:   User{},
		},
		{
			name: "插入失败,邮箱冲突",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx:     context.Background(),
			wantErr: ErrUserDuplicateEmail,
		},
		{
			name: "插入失败",
			sqlmock: func(t *testing.T) *sql.DB {
				db, mock, err := sqlmock.New()
				assert.NoError(t, err)
				mock.ExpectExec("INSERT INTO `users` .*").WillReturnError(errors.New("mock db error"))
				return db
			},
			ctx:     context.Background(),
			wantErr: errors.New("mock db error"),
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			sqldb := v.sqlmock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqldb,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			//初始化db不能出错 所以这里要断言必须为nil
			assert.NoError(t, err)
			dao := NewUserService(db)
			err = dao.Insert(v.ctx, v.u)
			assert.Equal(t, v.wantErr, err)
		})
	}
}
