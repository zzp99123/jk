package gormx

import (
	"github.com/prometheus/client_golang/prometheus"
	"gorm.io/gorm"
	"time"
)

type CallbackGorm struct {
	vector      *prometheus.SummaryVec
	Namespace   string
	Subsystem   string
	Names       string
	Help        string
	ConstLabels string
}

func NewCallbackGorm(Namespace string, Subsystem string, Names string, Help string, ConstLabels string) *CallbackGorm {
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: Namespace,
		Subsystem: Subsystem,
		Name:      Names,
		Help:      Help,
		ConstLabels: map[string]string{
			"db": ConstLabels,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.005,
			0.999: 0.0001,
		},
	},
		// 如果是 JOIN 查询，table 就是 JOIN 在一起的
		// 或者 table 就是主表，A JOIN B，记录的是 A
		[]string{"type", "table"})
	pcb := &CallbackGorm{
		vector: vector,
	}
	prometheus.MustRegister(vector)
	return pcb
}

// 计算增删改查运行时间
func (c *CallbackGorm) RegisterAll(db *gorm.DB) {
	err := db.Callback().Create().Before("*").Register("prometheus_create_before", c.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Create().After("*").Register("prometheus_create_after", c.After("creat"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().Before("*").Register("prometheus_update_before", c.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Update().After("*").Register("prometheus_update_after", c.After("update"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().Before("*").Register("prometheus_delete_before", c.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Delete().After("*").Register("prometheus_delete_delete", c.After("delete"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().Before("*").Register("prometheus_raw_before", c.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Raw().After("*").Register("prometheus_raw_after", c.After("raw"))
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().Before("*").Register("prometheus_row_before", c.Before())
	if err != nil {
		panic(err)
	}
	err = db.Callback().Row().After("*").Register("prometheus_row_after", c.After("row"))
	if err != nil {
		panic(err)
	}
}
func (c *CallbackGorm) Before() func(db *gorm.DB) {
	return func(db *gorm.DB) {
		startTime := time.Now()
		db.Set("start_time", startTime)
	}
}
func (c *CallbackGorm) After(typ string) func(db *gorm.DB) {
	return func(db *gorm.DB) {
		val, _ := db.Get("start_time")
		startTime, ok := val.(time.Time)
		if !ok {
			return
		}
		table := db.Statement.Table
		if table == "" {
			table = "unknow"
		}
		c.vector.WithLabelValues(typ, table).Observe(float64(time.Since(startTime).Milliseconds()))
	}
}

func (c *CallbackGorm) Name() string {
	return "prometheus-query"
}

func (c *CallbackGorm) Initialize(db *gorm.DB) error {
	c.RegisterAll(db)
	return nil
}
