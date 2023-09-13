package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"goFoundation/webook/internal/domain"
	"goFoundation/webook/internal/service"
	"net/http"
	"time"
)

var (
	userIdKey = "userId"
	bizLogin  = "login"
)

//var _ handler = &UserHandler{}
//
//// 这个更优雅
//var _ handler = (*UserHandler)(nil)

type UserHandler struct {
	svc         service.UserServiceIF
	svcCode     service.CodeServiceIF
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandle(svc service.UserServiceIF, svcCode service.CodeServiceIF) *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc,
		svcCode,
		emailExp,
		passwordExp,
	}

}
func (u *UserHandler) RegisterRouter(server *gin.Engine) {
	//up := server.Group("/user")
	server.POST("/users/signup", u.SignUp)
	//server.POST("/users/login", u.Login)
	server.POST("/users/login", u.LoginJwt)
	server.POST("/users/edit", u.Edit)
	server.GET("/users/profile", u.Profile)
	server.POST("/users/login_sms/code/send", u.SendLoginSMSCode)
	server.POST("/users/login_sms", u.LoginSMS)
}

// 发送
func (u *UserHandler) SendLoginSMSCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var r Req
	err := ctx.Bind(&r)
	if err != nil {
		return
	}
	if r.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "手机号为空",
		})
		return
	}
	err = u.svcCode.Set(ctx, bizLogin, r.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
}

// 验证码验证
func (u *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var r Req
	if err := ctx.Bind(&r); err != nil {
		return
	}
	ok, err := u.svcCode.Verify(ctx, bizLogin, r.Phone, r.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "验证码有误",
		})
		return
	}
	//你如果没注册过 直接先给你注册然后在登录 如果注册过直接登录
	res, err := u.svc.FindOrCreate(ctx, r.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	err = u.setJWTToken(ctx, res.Id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "验证码校验通过",
	})
}
func (u *UserHandler) SignUp(ctx *gin.Context) {
	//ctx.String(http.StatusOK, "123")
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var s SignUpReq
	//Bind 方法是 Gin 里面最常用的用于接收请求的方法。
	err := ctx.Bind(&s)
	if err != nil {
		return
	}
	//这种写法针对简单的正则表达式
	//ok, err := regexp.Match(emailRegexPattern, []byte(s.Email))
	ok, err := u.emailExp.MatchString(s.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "你的邮箱格式不对")
		return
	}
	if s.Password != s.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}
	ok, err = u.passwordExp.MatchString(s.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、特殊字符")
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    s.Email,
		Password: s.Password,
	})
	if err == service.ErrUserDuplicateEmail {
		ctx.String(http.StatusOK, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}
	//ctx.JSON(http.StatusOK, Result{
	//	Msg: "注册成功",
	//})
	ctx.String(http.StatusOK, "注册成功")
}

// 登录session登录方式
func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var l LoginReq
	err := ctx.Bind(&l)
	if err != nil {
		return
	}
	ok, err := u.svc.Login(ctx, l.Email, l.Password)
	if err == service.ErrInvaildUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//登录成功以后需要检验登录
	//先取session
	sess := sessions.Default(ctx)
	//存值
	sess.Set(userIdKey, ok.Id)
	//Gin Session 参数
	sess.Options(sessions.Options{
		MaxAge: 60, //负数就是退出登录 正数就是cookie失效时间
		//生产环境使用
		//Secure: true, cookie里面没有密码
		//HttpOnly: true, //标头不会携带cookie
	})
	//sess的机制必须执行
	sess.Save()
	ctx.String(http.StatusOK, "登录成功")
}

// jwt方式登录
func (u *UserHandler) LoginJwt(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var l LoginReq
	err := ctx.Bind(&l)
	if err != nil {
		return
	}
	ok, err := u.svc.Login(ctx, l.Email, l.Password)
	if err == service.ErrInvaildUserOrPassword {
		ctx.String(http.StatusOK, "用户名或密码不对")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	//创建一个精度512位的token
	if err = u.setJWTToken(ctx, ok.Id); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	fmt.Println("123", ok)
	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) setJWTToken(ctx *gin.Context, uid int64) error {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
		Uid:       uid,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		return err
	}
	ctx.Header("x-jwt-token", tokenStr)
	return nil
}

// 编辑
func (u *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Birthday        string `json:"birthday"`        //生日
		PersonalProfile string `json:"personalProfile"` //个人简介
		Nickname        string `json:"nickname"`        //昵称
	}
	var r Req
	err := ctx.Bind(&r)
	if err != nil {
		return
	}
	if r.Nickname == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "昵称不能为空",
		})
		return
	}
	if len(r.Nickname) < 3 || len(r.Nickname) > 20 {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "昵称必须在3到20个字符之间",
		})
		return
	}
	if len(r.PersonalProfile) >= 200 {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "个人简介的字数不能大于200个字符",
		})
		return
	}
	birthday, err := time.Parse(time.DateOnly, r.Birthday)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "日期格式不对",
		})
	}
	uc := ctx.MustGet("user").(UserClaims)
	err = u.svc.Edit(ctx, domain.User{
		Id:              uc.Uid,
		Nickname:        r.Nickname,
		Birthday:        birthday,
		PersonalProfile: r.PersonalProfile,
	})
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{Msg: "OK"})
}

// 查看
func (u *UserHandler) Profile(ctx *gin.Context) {
	type Req struct {
		Email string `json:"email"`
	}
	sess := sessions.Default(ctx)
	id := sess.Get(userIdKey).(int64)
	ok, err := u.svc.Profile(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//fmt.Println("1111111111111111111111")
	ctx.JSON(http.StatusOK, Req{
		Email: ok.Email,
	})
}

type UserClaims struct {
	jwt.RegisteredClaims
	// 声明你自己的要放进去 token 里面的数据
	Uid int64
	// 自己随便加
	UserAgent string
}
