package scheduler

import (
	"amazing_talker/pkg/model"
	"amazing_talker/pkg/model/ctype"
	"amazing_talker/pkg/model/option"
	"context"

	"github.com/rs/zerolog/log"
)

func (s *Scheduler) verifyIdentityAccountEmail(ctx context.Context) error {
	account, _, err := s.svc.ListIdentityAccounts(ctx, option.WhereIdentityAccountCondition{
		IdentityAccount: model.IdentityAccount{
			EmailVerifyStatus: ctype.VerifyInit,
		},
	})
	if err != nil {
		return err
	}

	for i := range account {
		err := s.svc.SendVerifyEmail(ctx, account[i])
		if err != nil {
			log.Ctx(ctx).Err(err).Msgf("Fail to send email account: %+v", account[i])
		}
	}

	return nil
}

func (s *Scheduler) verifyIdentityAccountPhone(ctx context.Context) error {
	account, _, err := s.svc.ListIdentityAccounts(ctx, option.WhereIdentityAccountCondition{
		IdentityAccount: model.IdentityAccount{
			PhoneVerifyStatus: ctype.VerifyInit,
		},
	})
	if err != nil {
		return err
	}

	for i := range account {
		err := s.svc.SendVerifyPhone(ctx, account[i])
		if err != nil {
			log.Ctx(ctx).Err(err).Msgf("Fail to send email account: %+v", account[i])
		}
	}

	return nil
}
