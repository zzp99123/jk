package startup

import "goFoundation/webook/pkg/logger"

func InitLog() logger.LoggerV1 {
	return &logger.NopLogger{}
}
