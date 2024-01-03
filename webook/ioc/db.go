package ioc

import (
	"fmt"
	"github.com/spf13/viper"
	dao2 "goFoundation/webook/interactive/repository/dao"
	"goFoundation/webook/internal/repository/dao"
	"goFoundation/webook/pkg/gormx"
	"goFoundation/webook/pkg/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
	"gorm.io/plugin/prometheus"
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
	//dsn := viper.GetString("db.mysql")
	//println(dsn)
	if err != nil {
		panic(fmt.Errorf("初始化配置失败 %v, 原因 %w", cfg, err))
	}
	//打印数据库日志
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
		// 使用 DEBUG 来打印
		// 缺了一个 writer
		//Logger: glogger.New(gormLoggerFunc(l.Debug), glogger.Config{
		//	// 慢查询阈值，只有执行时间超过这个阈值，才会使用
		//	// 50ms， 100ms
		//	// SQL 查询必然要求命中索引，最好就是走一次磁盘 IO
		//	// 一次磁盘 IO 是不到 10ms
		//	SlowThreshold:             time.Millisecond * 10,
		//	IgnoreRecordNotFoundError: true,
		//	ParameterizedQueries:      true,
		//	LogLevel:                  glogger.Info,
		//}),
	})
	if err != nil {
		// 我只会在初始化过程中 panic
		// panic 相当于整个 goroutine 结束
		// 一旦初始化过程出错，应用就不要启动了
		panic(err)
	}
	//统计 GORM 的执行时间
	err = db.Use(prometheus.New(prometheus.Config{
		DBName:          "webook",
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.MySQL{
				VariableNames: []string{"thread_running"},
			},
		},
	}))
	if err != nil {
		panic(err)
	}

	// 监控查询数据库增删改查的执行时间
	pcb := gormx.NewCallbackGorm("geekbang_daming", "webook", "gorm_query_time", "统计 GORM 的执行时间", "webook")
	db.Use(pcb)
	//在 GORM 中接入 OpenTelemetry
	db.Use(tracing.NewPlugin(tracing.WithDBName("webook"),
		tracing.WithQueryFormatter(func(query string) string {
			l.Debug("", logger.String("query", query))
			return query

		}),
		// 不要记录 metrics
		tracing.WithoutMetrics(),
		// 不要记录查询参数
		tracing.WithoutQueryVariables()))

	//dao.NewUserDAOV1(func() *gorm.DB {
	//viper.OnConfigChange(func(in fsnotify.Event) {
	//oldDB := db
	//db, err = gorm.Open(mysql.Open())
	//pt := unsafe.Pointer(&db)
	//atomic.StorePointer(&pt, unsafe.Pointer(&db))
	//oldDB.Close()
	//})
	// 要用原子操作
	//return db
	//})

	err = dao.InitTable(db)
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
