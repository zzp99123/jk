package wechat

import (
	"context"
	"goFoundation/webook/internal/domain"
)

type Service interface {
	AuthUrl(ctx context.Context, state string) (string, error)
	// VerifyCode 校验扫码登录完成的url 这个code你可以校验unionId也可以是openId
	// 目前大部分公司的 OAuth2 平台都差不多的设计
	// 返回一个 unionId。这个你可以理解为，在第三方平台上的 unionId
	// 你也可以考虑使用 openId 来替换
	// 一家公司如果有很多应用，不同应用都有自建的用户系统
	// 那么 openId 可能更加合适
	VerifyCode(ctx context.Context, code string) (domain.WechatInfo, error)
}
