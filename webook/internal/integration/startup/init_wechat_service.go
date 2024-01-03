package startup

import (
	"goFoundation/webook/internal/service/oauth2/wechat"
	"goFoundation/webook/pkg/logger"
	"os"
)

func InitPhantomWechatService(l logger.LoggerV1) wechat.Service {
	//获取指定键对应的环境变量值
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量WECHAT_APP_ID")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET ")
	}
	return wechat.NewServiceWechat(appId, appKey, l)
}
