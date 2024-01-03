// 微信扫码登录
package web

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"goFoundation/webook/internal/service"
	"goFoundation/webook/internal/service/oauth2/wechat"
	ijwt "goFoundation/webook/internal/web/jwt"
	"net/http"
	"time"
)

type OAuth2WechatHandler struct {
	wechatService wechat.Service
	userService   service.UserService
	ijwt.Handler
	stateKey []byte
	cfg      WechatHandlerConfig
}
type WechatHandlerConfig struct {
	Secure bool
	//StateKey
}

func NewOAuth2WechatHandler(wechatService wechat.Service, userService service.UserService, cfg WechatHandlerConfig, jwtHdl ijwt.Handler) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		wechatService: wechatService,
		userService:   userService,
		stateKey:      []byte("95osj3fUD7fo1mlYdDbncXz4VD2igvf0"),
		cfg:           cfg,
		Handler:       jwtHdl,
	}
}
func (o *OAuth2WechatHandler) RegisterRouters(s *gin.Engine) {
	up := s.Group("/oauth2/wechat")
	up.GET("/authurl", o.authurl)
	up.Any("/callback", o.callback)
}

// 扫码登录
func (o *OAuth2WechatHandler) authurl(ctx *gin.Context) {
	state := uuid.New()
	res, err := o.wechatService.AuthUrl(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "构造扫码登录URL失败",
		})
		return
	}
	if err := o.setStateCookie(ctx, state); err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统异常",
		})
		return
	}
	//返回地址
	ctx.JSON(http.StatusOK, Result{
		Data: res,
	})
	return
}

// 校验
func (o *OAuth2WechatHandler) callback(ctx *gin.Context) {
	//code是我Query读到的
	code := ctx.Query("code")
	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "登录失败",
		})
		return
	}
	//校验通过可以拿到那2个id
	info, err := o.wechatService.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//根据 OpenId 和 UnionId来判断你这个扫码完成以后校验 然后在去数据库里看你注没注册过
	//这就可以拿到id
	res, err := o.userService.FindOrCreateByWechat(ctx, info)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = o.SetLoginToken(ctx, res.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
}

// 存state
func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, StateClaims{
		State: state,
		//过期时间
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 10)),
		},
	})
	tokenStr, err := token.SignedString(o.stateKey)
	if err != nil {
		return err
	}
	//存到cookie中
	ctx.SetCookie("jwt-state", tokenStr,
		600, "/oauth2/wechat/callback",
		"", o.cfg.Secure, true)
	return nil
}
func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie("jwt-state")
	if err != nil {
		return fmt.Errorf("拿不到 state 的 cookie, %w", err)
	}
	var s StateClaims
	//这就相当于我存起来加密后的statekey经过处理给到StateClaims里面state 然后在拿着StateClaims里面state和我在ctx.Query的state做比较
	token, err := jwt.ParseWithClaims(ck, &s, func(token *jwt.Token) (interface{}, error) {
		return o.stateKey, nil
	})
	if err != nil || !token.Valid {
		return fmt.Errorf("token 已经过期了, %w", err)
	}
	if s.State != state {
		return errors.New("state 不相等")
	}
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string
}

//校验state
