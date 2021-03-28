package service

import (
	"amazing_talker/internal/errors"
	"amazing_talker/pkg/model"
	"amazing_talker/pkg/model/ctype"
	"amazing_talker/pkg/model/option"
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"gitlab.com/howmay/gopher/db"
	"golang.org/x/text/language"
	"google.golang.org/api/gmail/v1"
	"gorm.io/gorm"
)

var (
	chineseI18n = i18n.Message{
		ID:    "Emails",
		One:   "{{.Name}} 您好! 歡迎加入 AmazingTalker",
		Other: "{{.Name}} 您好! 歡迎加入 AmazingTalker",
	}

	englishI18n = i18n.Message{
		ID:    "Emails",
		One:   "Hi {{.Name}} ! Welcome to AmazingTalker.",
		Other: "Hi {{.Name}} ! Welcome to AmazingTalker.",
	}

	verifyEmailTemplate = "From: amazingTalker@gmail.com\r\n" +
		"To: %s \r\n" +
		"Subject: =?utf-8?B?%s?=\r\n\r\n" +
		"%s"
)

func (s *service) getEmailBodyMessage(lang string, name string) string {
	loc := i18n.NewLocalizer(s.bundle, language.Chinese.String())

	if strings.ContainsAny(lang, "en") {
		loc = i18n.NewLocalizer(s.bundle, language.English.String())
	}

	translation := loc.MustLocalize(&i18n.LocalizeConfig{
		MessageID: "Emails",
		TemplateData: map[string]interface{}{
			"Name": name,
		},
	})

	return translation
}

func (s *service) SendVerifyEmail(ctx context.Context, account model.IdentityAccount) error {
	bodyMessage := s.getEmailBodyMessage(account.AcceptLanguage, account.Name)

	txErr := db.ExecuteTx(ctx, s.repo.WriteDB(), func(tx *gorm.DB) error {
		if err := s.UpdateIdentityAccount(ctx, option.UpdateIdentityAccountCondition{
			WhereCondition: option.WhereIdentityAccountCondition{
				IdentityAccount: model.IdentityAccount{
					ID:                account.ID,
					EmailVerifyStatus: ctype.VerifyInit,
				},
			},
			UpdateColumn: option.UpdateIdentityAccountColumn{
				EmailVerifyStatus: ctype.VerifySend,
				SendEmailVerifyAt: time.Now().UTC(),
			},
		}); err != nil {
			return err
		}

		if err := s.SendEmail(account.Email, bodyMessage, bodyMessage); err != nil {
			return errors.NewWithMessagef(errors.ErrInternalError, "fail to send mail to [ %s ], err: %+v", account.Email, err)
		}
		return nil
	})
	return txErr
}

// SendEmail ...
func (s *service) SendEmail(toEmail, subject, body string) error {
	messageStr := []byte(fmt.Sprintf(verifyEmailTemplate, toEmail, base64.StdEncoding.EncodeToString([]byte(subject)), body))

	msg := gmail.Message{
		Raw: base64.StdEncoding.EncodeToString(messageStr),
	}

	_, err := s.gmailSvc.Users.Messages.Send("me", &msg).Do()
	if err != nil {
		return errors.NewWithMessagef(errors.ErrInternalError, "fail to send mail to [ %s ], err: %+v", toEmail, err)
	}

	return err
}
