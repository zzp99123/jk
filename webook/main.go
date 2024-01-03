package main

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"goFoundation/webook/ioc"
	"goFoundation/webook/pkg/viperx"
	"net/http"
	"time"
)

func main() {
	viperx.InitViperV1()
	//initLogger()
	initPrometheus()
	//在 Gin 中接入 OpenTelemetry
	closeFunc := ioc.InitOTEL()
	//AllKeys 返回包含值的所有键，无论它们设置在何处
	//keys := viper.AllKeys()
	//println(keys)
	//您可能需要将 viper 中保存的所有设置编组到字符串中，而不是将它们写入文件中
	//setting := viper.AllSettings()
	//fmt.Println(setting)
	app := InitWebServer()
	app.cron.Start()
	// Consumer 在我设计下，类似于 Web，或者 GRPC 之类的，是一个顶级入口
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}
	server := app.web
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "你好，你来了")
	})

	server.Run(":8080")
	// 作业
	//server.Run(":8081")
	// 一分钟内你要关完，要退出
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	closeFunc(ctx)
	c := app.cron.Stop()
	//超时退出
	k := time.NewTicker(time.Minute * 10)
	select {
	case <-k.C:
	case <-c.Done():

	}
}
func initPrometheus() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(":8081", nil)
	}()
}
func initLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	zap.L().Info("这是 replace 之前")
	// 如果你不 replace，直接用 zap.L()，你啥都打不出来。
	zap.ReplaceGlobals(logger)
	zap.L().Info("hello，你搞好了")

	type Demo struct {
		Name string `json:"name"`
	}
	zap.L().Info("这是实验参数",
		zap.Error(errors.New("这是一个 error")),
		zap.Int64("id", 123),
		zap.Any("一个结构体", Demo{Name: "hello"}))
}
