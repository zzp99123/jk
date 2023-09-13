package middleware

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type LoginJwtMiddleBuild struct {
	path []string
}

func NewUserJwtService() *LoginJwtMiddleBuild {
	return &LoginJwtMiddleBuild{}
}
func (l *LoginJwtMiddleBuild) LoginPath(path string) *LoginJwtMiddleBuild {
	l.path = append(l.path, path)
	return l
}

func (l *LoginJwtMiddleBuild) Build() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		for _, v := range l.path {
			if ctx.Request.URL.Path == v {
				return
			}
		}
		//jwt检验
		//先获取这个header
		tokenHeader := ctx.GetHeader("Authorization")
		//没带Authorization登录
		if tokenHeader == "" {
			//没带Authorization登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//fmt.Println("11111111111111111111", tokenHeader)
		//如果登录过的花 先切割
		segs := strings.Split(tokenHeader, " ")
		if len(segs) > 2 {
			//格式不对登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenstr := segs[1]
		//fmt.Println("222222222222222222222222", tokenstr)
		token, err := jwt.Parse(tokenstr, func(token *jwt.Token) (interface{}, error) {

			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})
		//fmt.Println("3333333333333333333333", token)
		if err != nil {
			//没人登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid {
			//没人登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

	}
}
