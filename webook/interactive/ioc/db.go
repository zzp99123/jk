package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	dao2 "goFoundation/webook/interactive/repository/dao"
	"goFoundation/webook/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(l logger.LoggerV1) *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	//默认配置
	var cfg = Config{
		DSN: "root:root@tcp(localhost:13316)/webook_default",
	}
	// 看起来，remote 不支持 key 的切割
	err := viper.UnmarshalKey("db", &cfg)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", cfg, err))
	}
	//打印数据库日志
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao2.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}

type gormLoggerFunc func(msg string, fields ...logger.Field)

func (g gormLoggerFunc) Printf(msg string, args ...interface{}) {
	g(msg, logger.Field{Key: "args", Value: args})
}

type DoSomething interface {
	DoABC() string
}

type DoSomethingFunc func() string

func (d DoSomethingFunc) DoABC() string {
	return d()
}
