package loggerxx

import "go.uber.org/zap"

var Logger *zap.Logger

func InitLogger(l *zap.Logger) {
	Logger = l
}
func InitLoggerV1() {
	Logger, _ = zap.NewDevelopment()
}
