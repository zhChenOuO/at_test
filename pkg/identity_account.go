package pkg

import (
	"amazing_talker/pkg/model"
	option "amazing_talker/pkg/model/option"
	"context"

	gorm "gorm.io/gorm"
)

// IdentityAccountService service介面層
type IdentityAccountService interface {
	GetIdentityAccount(ctx context.Context, opt option.WhereIdentityAccountCondition) (model.IdentityAccount, error)
	CreateIdentityAccount(ctx context.Context, data *model.IdentityAccount) error
	ListIdentityAccounts(ctx context.Context, opt option.WhereIdentityAccountCondition) ([]model.IdentityAccount, int64, error)
	UpdateIdentityAccount(ctx context.Context, opt option.UpdateIdentityAccountCondition) error
	DeleteIdentityAccount(ctx context.Context, opt option.WhereIdentityAccountCondition) error
}

// IdentityAccountRepository repository介面層
type IdentityAccountRepository interface {
	GetIdentityAccount(ctx context.Context, tx *gorm.DB, opt option.WhereIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) (model.IdentityAccount, error)
	CreateIdentityAccount(ctx context.Context, tx *gorm.DB, data *model.IdentityAccount, scopes ...func(*gorm.DB) *gorm.DB) error
	ListIdentityAccounts(ctx context.Context, tx *gorm.DB, opt option.WhereIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) ([]model.IdentityAccount, int64, error)
	UpdateIdentityAccount(ctx context.Context, tx *gorm.DB, opt option.UpdateIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) error
	DeleteIdentityAccount(ctx context.Context, tx *gorm.DB, opt option.WhereIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) error
}
