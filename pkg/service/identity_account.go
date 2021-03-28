package service

import (
	"amazing_talker/internal/errors"
	model "amazing_talker/pkg/model"
	option "amazing_talker/pkg/model/option"
	"context"

	"gitlab.com/howmay/gopher/db"
	"gorm.io/gorm"
)

// GetIdentityAccount 取得IdentityAccount的資訊
func (s *service) GetIdentityAccount(ctx context.Context, opt option.WhereIdentityAccountCondition) (model.IdentityAccount, error) {
	return s.repo.GetIdentityAccount(ctx, nil, opt)
}

// CreateIdentityAccount 建立IdentityAccount
func (s *service) CreateIdentityAccount(ctx context.Context, data *model.IdentityAccount) error {
	txErr := db.ExecuteTx(ctx, s.repo.WriteDB(), func(tx *gorm.DB) error {
		switch {
		case data.Email != "":
			_, err := s.repo.GetIdentityAccount(ctx, tx, option.WhereIdentityAccountCondition{
				IdentityAccount: model.IdentityAccount{
					Email: data.Email,
				},
			})
			if err == nil {
				return errors.WithStack(errors.ErrEmailAlreadyExists)
			} else if !errors.Is(err, errors.ErrResourceNotFound) {
				return err
			}
		case data.Phone != "":
			_, err := s.repo.GetIdentityAccount(ctx, tx, option.WhereIdentityAccountCondition{
				IdentityAccount: model.IdentityAccount{
					PhoneAreaCode: data.PhoneAreaCode,
					Phone:         data.Phone,
				},
			})
			if err == nil {
				return errors.WithStack(errors.ErrPhoneAlreadyExists)
			} else if !errors.Is(err, errors.ErrResourceNotFound) {
				return err
			}
		}

		if err := s.repo.CreateIdentityAccount(ctx, tx, data); err != nil {
			return err
		}
		return nil
	})

	return txErr
}

// ListIdentityAccounts 列出IdentityAccount
func (s *service) ListIdentityAccounts(ctx context.Context, opt option.WhereIdentityAccountCondition) ([]model.IdentityAccount, int64, error) {
	return s.repo.ListIdentityAccounts(ctx, nil, opt)
}

// UpdateIdentityAccount 更新IdentityAccount
func (s *service) UpdateIdentityAccount(ctx context.Context, opt option.UpdateIdentityAccountCondition) error {
	return s.repo.UpdateIdentityAccount(ctx, nil, opt)
}

// DeleteIdentityAccount 刪除IdentityAccount
func (s *service) DeleteIdentityAccount(ctx context.Context, opt option.WhereIdentityAccountCondition) error {
	return s.repo.DeleteIdentityAccount(ctx, nil, opt)
}
