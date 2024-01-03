package wechat

import (
	"context"
	"encoding/json"
	"fmt"
	"goFoundation/webook/pkg/logger"

	"goFoundation/webook/internal/domain"
	"net/http"
	"net/url"
)

const authURLPattern = "https://open.weixin.qq.com/connect/qrconnect?appid=%s&redirect_uri=%s&response_type=code&scope=snsapi_login&state=%s#wechat_redire"

var redirectURL = url.PathEscape("https://meoying.com/oauth2/wechat/callback")

type serviceWechat struct {
	appId     string
	appSecret string
	client    *http.Client
	l         logger.LoggerV1
}

func NewServiceWechat(appId string, appSecret string, l logger.LoggerV1) Service {
	return &serviceWechat{
		appId:     appId,
		appSecret: appSecret,
		client:    http.DefaultClient,
		l:         l,
	}
}
func (s *serviceWechat) AuthUrl(ctx context.Context, state string) (string, error) {
	return fmt.Sprintf(authURLPattern, s.appId, redirectURL, state), nil
}

func (s *serviceWechat) VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error) {
	//先发请求
	const urlPatten = "https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code"
	target := fmt.Sprintf(urlPatten, s.appId, s.appSecret, code)
	//可以用get发请求 也可用
	res, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	//发送请求
	resp, err := s.client.Do(res)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	//判断resq.body里面的值
	//json.NewDecoder用于http连接与socket连接的读取与写入，或者文件读取
	//2、json.Unmarshal用于直接是byte的输入
	decoder := json.NewDecoder(resp.Body)
	var r Result
	err = decoder.Decode(&r)
	if err != nil {
		return domain.WechatInfo{}, err
	}
	if r.ErrCode != 0 {
		return domain.WechatInfo{}, fmt.Errorf("微信返回错误响应，错误码%d，错误信息%s", r.ErrCode, r.ErrMsg)
	}
	return domain.WechatInfo{
		OpenId:  r.OpenId,
		UnionId: r.UnionId,
	}, nil
}

type Result struct {
	ErrCode int64  `json:"errcode"`
	ErrMsg  string `json:"errMsg"`

	Scope string `json:"scope"`

	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`

	OpenId  string `json:"openid"`
	UnionId string `json:"unionid"`
}

//正确返回这个
//{
//"access_token":"ACCESS_TOKEN",
//"expires_in":7200,
//"refresh_token":"REFRESH_TOKEN",
//"openid":"OPENID",
//"scope":"SCOPE",
//"unionid": "UNIONID"
//}

//错误返回这个
//{
//"errcode":40029,"errmsg":"invalid code"
//}
