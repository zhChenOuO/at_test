package repository

import (
	"amazing_talker/internal/errors"
	"amazing_talker/pkg/model"
	"amazing_talker/pkg/model/option"
	"context"
	"reflect"

	"gitlab.com/howmay/gopher/db"
	"gorm.io/gorm"
)

// GetIdentityAccount 取得IdentityAccount的資訊
func (repo *repository) GetIdentityAccount(ctx context.Context, tx *gorm.DB, opt option.WhereIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) (model.IdentityAccount, error) {
	if tx == nil {
		tx = repo.readDB
	}
	tx = tx.Scopes(scopes...)
	var IdentityAccount model.IdentityAccount
	err := tx.Table(IdentityAccount.TableName()).Scopes(opt.Where).First(&IdentityAccount).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return IdentityAccount, errors.WithStack(errors.ErrResourceNotFound)
		}
		return IdentityAccount, errors.Wrapf(errors.ErrInternalError, "database: IdentityAccount err: %s", err.Error())
	}
	return IdentityAccount, nil
}

// CreateIdentityAccount 建立IdentityAccount
func (repo *repository) CreateIdentityAccount(ctx context.Context, tx *gorm.DB, data *model.IdentityAccount, scopes ...func(*gorm.DB) *gorm.DB) error {
	if tx == nil {
		tx = repo.writeDB
	}

	data.Init()

	tx = tx.Scopes(scopes...)
	err := tx.Create(data).Error
	if err != nil {
		if db.IsDuplicateErr(err) {
			return errors.Wrapf(errors.ErrResourceAlreadyExists, "database: CreateIdentityAccount err: %s", err.Error())
		}
		return errors.Wrapf(errors.ErrInternalError, "database: CreateIdentityAccount err: %s", err.Error())
	}
	return nil
}

// ListIdentityAccounts 列出IdentityAccount
func (repo *repository) ListIdentityAccounts(ctx context.Context, tx *gorm.DB, opt option.WhereIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) ([]model.IdentityAccount, int64, error) {
	if tx == nil {
		tx = repo.readDB
	}
	tx = tx.Scopes(scopes...)
	var IdentityAccounts []model.IdentityAccount
	var total int64
	db := tx.Table(model.IdentityAccount{}.TableName()).Scopes(opt.Where)
	err := db.Count(&total).Error
	if err != nil {
		return nil, total, errors.Wrapf(errors.ErrInternalError, "database: ListIdentityAccount err: %s", err.Error())
	}
	err = db.Scopes(opt.Pagination.LimitAndOffset, opt.Sorting.Sort).Find(&IdentityAccounts).Error
	if err != nil {
		return nil, total, errors.Wrapf(errors.ErrInternalError, "database: ListIdentityAccount err: %s", err.Error())
	}
	return IdentityAccounts, total, nil
}

// UpdateIdentityAccount 更新IdentityAccount
func (repo *repository) UpdateIdentityAccount(ctx context.Context, tx *gorm.DB, opt option.UpdateIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) error {
	if tx == nil {
		tx = repo.writeDB
	}
	tx = tx.Scopes(scopes...)
	if reflect.DeepEqual(opt.WhereCondition, option.WhereIdentityAccountCondition{}) {
		return errors.Wrap(errors.ErrInternalError, "database: UpdateIdentityAccount err: where condition can't empty")
	}
	err := tx.Table(model.IdentityAccount{}.TableName()).Scopes(opt.WhereCondition.Where).Updates(opt.UpdateColumn).Error
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "database: UpdateIdentityAccount err: %s", err.Error())
	}
	return nil
}

// DeleteIdentityAccount 刪除IdentityAccount
func (repo *repository) DeleteIdentityAccount(ctx context.Context, tx *gorm.DB, opt option.WhereIdentityAccountCondition, scopes ...func(*gorm.DB) *gorm.DB) error {
	if tx == nil {
		tx = repo.writeDB
	}
	tx = tx.Scopes(scopes...)
	if reflect.DeepEqual(opt.IdentityAccount, model.IdentityAccount{}) {
		return errors.Wrap(errors.ErrInvalidInput, "database: DeleteIdentityAccount err: WhereIdentityAccountCondition is empty")
	}
	err := tx.Scopes(opt.Where).Delete(&model.IdentityAccount{}).Error
	if err != nil {
		return errors.Wrapf(errors.ErrInternalError, "database: DeleteIdentityAccount err: %s", err.Error())
	}
	return nil
}
