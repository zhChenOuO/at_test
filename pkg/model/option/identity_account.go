package option

import (
	model "amazing_talker/pkg/model"
	"amazing_talker/pkg/model/ctype"
	"time"

	"gitlab.com/howmay/gopher/common"
	"gorm.io/gorm"
)

// WhereIdentityAccountCondition ORM查詢條件
type WhereIdentityAccountCondition struct {
	IdentityAccount model.IdentityAccount `json:"identity_account"`
	Pagination      common.Pagination     `json:"pagination"`
	BaseWhere       common.BaseWhere      `json:"base_where"`
	Sorting         common.Sorting        `json:"sorting"`
}

// Where 基礎的查詢條件
func (where *WhereIdentityAccountCondition) Where(db *gorm.DB) *gorm.DB {
	db = db.Where(where.IdentityAccount)
	db = db.Scopes(where.BaseWhere.Where)
	return db
}

// UpdateIdentityAccountCondition ORM更新條件
type UpdateIdentityAccountCondition struct {
	WhereCondition WhereIdentityAccountCondition
	UpdateColumn   UpdateIdentityAccountColumn
}

// UpdateIdentityAccountColumn ORM更新欄位
type UpdateIdentityAccountColumn struct {
	PhoneVerifyStatus ctype.VerifyStatus `json:"phone_verify_status"`
	SendPhoneVerifyAt time.Time          `json:"send_phone_verify_at"`
	EmailVerifyStatus ctype.VerifyStatus `json:"email_verify_status"`
	SendEmailVerifyAt time.Time          `json:"send_email_verify_at"`
}
