// 根据情况不同简易封装的zap的Write方法
package zapx

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MyCore struct {
	zapcore.Core
}

func (c MyCore) Write(entry zapcore.Entry, fds []zapcore.Field) error {
	for _, fd := range fds {
		//打日志遇到敏感类似手机号不能全部暴露 176****2527
		if fd.Key == "phone" {
			phone := fd.String
			fd.String = phone[:3] + "****" + phone[7:]
		}
	}
	return c.Core.Write(entry, fds)
}

func MaskPhone(key string, value string) zap.Field {
	value = value[:3] + "****" + value[7:]
	return zap.Field{
		Key:    key,
		String: value,
	}
}
