package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"goFoundation/webook/internal/service/sms"
)

type ServiceAuth struct {
	svc sms.Service
	key string
}
type Claims struct {
	jwt.RegisteredClaims
	tpl string //真正要用的token
}

func (s *ServiceAuth) Send(ctx context.Context, biz string, args []string, numbers ...string) error {
	var c Claims
	tocker, err := jwt.ParseWithClaims(biz, &c, func(token *jwt.Token) (interface{}, error) {
		return s.key, nil
	})
	if err != nil {
		return err
	}
	if !tocker.Valid {
		return errors.New("tocken 不合适")
	}
	return s.svc.Send(ctx, c.tpl, args, numbers...)
}
