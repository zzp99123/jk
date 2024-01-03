package ioc

import (
	"goFoundation/webook/internal/service/sms"
	"goFoundation/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	//接入监控
	//return metrics.NewPrometheusDecorator(memory.NewMemoryService())
	return memory.NewMemoryService()
}
