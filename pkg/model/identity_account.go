package model

import (
	"amazing_talker/configuration"
	"amazing_talker/internal/claims"
	"amazing_talker/internal/errors"
	"amazing_talker/pkg/model/ctype"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gitlab.com/howmay/gopher/db"
)

// IdentityAccount ...
type IdentityAccount struct {
	ID                int                `json:"id"`
	Name              string             `json:"name"`
	Email             string             `json:"email"`
	Phone             string             `json:"phone"`
	PhoneAreaCode     string             `json:"phone_area_code"`
	Password          db.Crypto          `json:"password"`
	PhoneVerifyStatus ctype.VerifyStatus `json:"phone_verify_status"`
	SendPhoneVerifyAt time.Time          `json:"send_phone_verify_at"`
	EmailVerifyStatus ctype.VerifyStatus `json:"email_verify_status"`
	SendEmailVerifyAt time.Time          `json:"send_email_verify_at"`
	AcceptLanguage    string             `json:"accept_language"`
	CreatedAt         time.Time          `json:"created_at"`
	UpdatedAt         time.Time          `json:"updated_at"`
}

// TableName ...
func (IdentityAccount) TableName() string {
	return "identity_accounts"
}

// CreateToken ...
func (identity *IdentityAccount) CreateToken(cfg *configuration.App) (string, error) {
	if cfg.JwtExpireSec == 0 {
		cfg.JwtExpireSec = 3600 * 24
	}

	claims := &claims.Claims{
		ID:    identity.ID,
		Name:  identity.Name,
		Email: identity.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(time.Duration(cfg.JwtExpireSec) * time.Second).Unix(),
		},
	}

	jwtClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := jwtClaims.SignedString([]byte(cfg.JwtSecrets))
	if err != nil {
		return "", errors.WithMessagef(errors.ErrInternalError, "err: %s", err.Error())
	}
	return t, nil
}

// Init 初始化資料
func (identity *IdentityAccount) Init() {
	identity.PhoneVerifyStatus = ctype.VerifyInit
	identity.EmailVerifyStatus = ctype.VerifyInit
}
