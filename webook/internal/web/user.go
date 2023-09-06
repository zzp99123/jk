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

type UserHandler struct {
	svc         *service.ServiceUser
	emaileExp   *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandle(svc *service.ServiceUser) *UserHandler {
	const (
		emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
		// 和上面比起来，用 ` 看起来就比较清爽
		passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	emaileExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(passwordRegexPattern, regexp.None)
	return &UserHandler{
		svc,
		emaileExp,
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
	ok, err := u.emaileExp.MatchString(s.Email)
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
			Msg:  "邮箱不正确",
		})
		return
	}
	if s.Password != s.ConfirmPassword {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "密码不一致",
		})
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
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "密码不正确",
		})
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    s.Email,
		Password: s.Password,
	})
	if err == service.ErruserDuplicateEmail {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "重复邮箱，请换个邮箱",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "服务器异常，注册失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Msg: "注册成功",
	})
	fmt.Printf("%v", s)
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
	if err == service.ErrInvaildUserOrPasswod {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "用户名或密码不正确",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
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
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
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
	if err == service.ErrInvaildUserOrPasswod {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "用户名或密码不正确",
		})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		return
	}
	//创建一个精度512位的token
	if err = u.setJWTToken(ctx, ok.Id); err != nil {
		return
	}
	fmt.Println("123", ok)
	ctx.JSON(http.StatusOK, Result{
		Msg: "登录成功",
	})
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