package repository

import (
	"context"
	"goFoundation/webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany        = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = cache.ErrCodeVerifyTooManyTimes
)

type CodeRepositoryIF interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error)
}
type codeRepository struct {
	c cache.CodeCacheIF
}

func NewRepositoryCode(c cache.CodeCacheIF) CodeRepositoryIF {
	return &codeRepository{
		c: c,
	}
}
func (r *codeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return r.c.Set(ctx, biz, phone, code)
}
func (r *codeRepository) Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error) {
	return r.c.Verify(ctx, biz, phone, expectedCode)
}
