package main

import (
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"goFoundation/webook/pkg/saramax"
)

type App struct {
	web       *gin.Engine
	consumers []saramax.ConsumerEvents
	cron      *cron.Cron
}
