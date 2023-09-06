package middleware

import (
	"encoding/gob"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddleBuild struct {
	path []string
}

func NewUserService() *LoginMiddleBuild {
	return &LoginMiddleBuild{}
}

// 路径方法
func (l *LoginMiddleBuild) LoginPath(path string) *LoginMiddleBuild {
	l.path = append(l.path, path)
	return l
}
func (l *LoginMiddleBuild) Build() gin.HandlerFunc {
	gob.Register(time.Now()) //一旦你处理接口，你必须首先注册你的具体类型。
	return func(ctx *gin.Context) {
		//登录和注册不需要校验
		//if ctx.Request.URL.Path == "/users/signup" || ctx.Request.URL.Path == "/users/login" {
		//	return
		//}
		for _, v := range l.path {
			if ctx.Request.URL.Path == v {
				return
			}
		}
		sess := sessions.Default(ctx)
		//已经有sessions的插件了所以这个没用
		//if sess == nil {
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		id := sess.Get("userId")
		if id == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//刷新登录状态
		//村的时间
		updateTime := sess.Get("updatetime")
		sess.Set("userId", id)
		sess.Options(sessions.Options{
			MaxAge: 60,
		})
		now := time.Now().UnixMilli() //当前时间
		//第一次登录 还没刷新过
		if updateTime == nil {
			sess.Set("updatetime", now)
			err := sess.Save()
			if err != nil {
				panic(err)
			}
		}
		//已经登录过
		//类型断言 先确定类型
		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if now-updateTimeVal > 60*100 { //当前时间-存起来的时间 大于1分钟 我就刷新
			sess.Set("updatetime", now)
			err := sess.Save()
			if err != nil {
				panic(err)
			}
		}

	}
}
