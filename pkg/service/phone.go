package service

import (
	"amazing_talker/pkg/model"
	"context"
)

func (s *service) SendVerifyPhone(ctx context.Context, account model.IdentityAccount) error {
	// TODO 尚未實作 sms 驗證 
	return nil
}
