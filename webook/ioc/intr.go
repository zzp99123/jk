package ioc

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	intrv1 "goFoundation/webook/api/proto/gen/intr/v1"
	"goFoundation/webook/interactive/service"
	"goFoundation/webook/internal/web/client"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrGRPCClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type config struct {
		addr      string
		Secure    bool //表示创建的 Cookie 会被以安全的形式向服务器传输
		Threshold int32
	}
	var c config
	err := viper.UnmarshalKey("grpc.client.intr", &c)
	if err != nil {
		panic(err)
	}
	var d []grpc.DialOption
	if c.Secure {
		// 上面，要去加载你的证书之类的东西
		// 启用 HTTPS
	} else {
		d = append(d, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial("grpc-clinet", d...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewIntrLocalClient(svc)
	res := client.NewGreyScaleIntrClinet(remote, local)
	//在这里进行监听
	viper.OnConfigChange(func(in fsnotify.Event) {
		var c config
		err = viper.UnmarshalKey("grpc.client.intr", &c)
		if err != nil {
			// 你可以输出日志
		}
		res.UpdateThreshold(c.Threshold)
	})
	return res
}
