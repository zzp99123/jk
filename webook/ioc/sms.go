package ioc

import (
	"goFoundation/webook/internal/service/sms"
	"goFoundation/webook/internal/service/sms/memory"
)

func InitSMSService() sms.Service {
	return memory.NewMemoryService()
}
