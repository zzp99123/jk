package service

import (
	"context"
	"fmt"
	"goFoundation/webook/internal/repository"
	"goFoundation/webook/internal/service/sms"
	"math/rand"
)

const codeTplId = "1877556"

var (
	ErrCodeSendTooMany        = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooManyTimes = repository.ErrCodeVerifyTooManyTimes
)

type CodeService interface {
	Set(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error)
}
type codeService struct {
	r      repository.CodeRepository
	smsSvc sms.Service
}

func NewServiceCode(r repository.CodeRepository, smsSvc sms.Service) CodeService {
	return &codeService{
		r:      r,
		smsSvc: smsSvc,
	}
}

// 发送
func (s *codeService) Set(ctx context.Context, biz, phone string) error {
	code := s.Code()
	err := s.r.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	//这相当于验证码生产成功了 下一步需要发送
	err = s.smsSvc.Send(ctx, codeTplId, []string{code}, phone) //这步就是发送了
	// 这个地方怎么办？
	// 这意味着，Redis 有这个验证码，但是不好意思，
	// 我能不能删掉这个验证码？
	// 你这个 err 可能是超时的 err，你都不知道，发出了没
	// 在这里重试
	// 要重试的话，初始化的时候，传入一个自己就会重试的 smsSvc
	return err
}

// 接受验证码并验证
func (s *codeService) Verify(ctx context.Context, biz, phone, expectedCode string) (bool, error) {
	return s.r.Verify(ctx, biz, phone, expectedCode)
}

// 创建验证码
func (s *codeService) Code() string {
	// 0-999999
	num := rand.Intn(1000000)
	return fmt.Sprintf("%06d", num)
}
