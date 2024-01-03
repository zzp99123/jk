package wechat

import (
	"context"
	"github.com/stretchr/testify/require"
	"goFoundation/webook/pkg/logger"
	"os"
	"testing"
)

func Test_wecaht(t *testing.T) {
	appId, ok := os.LookupEnv("WECHAT_APP_ID")
	if !ok {
		panic("没有找到环境变量WECHAT_APP_ID")
	}
	appKey, ok := os.LookupEnv("WECHAT_APP_SECRET")
	if !ok {
		panic("没有找到环境变量 WECHAT_APP_SECRET ")
	}
	svc := NewServiceWechat(appId, appKey, &logger.NopLogger{})
	res, err := svc.VerifyCode(context.Background(), "051D6b000Yn4FQ14Rd300FgOF33D6b0s")
	require.NoError(t, err)
	t.Log(res)
}
