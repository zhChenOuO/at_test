package service

import (
	pkg "amazing_talker/pkg"
	model "amazing_talker/pkg/model"
	"context"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/api/gmail/v1"
)

func Test_service_CreateIdentityAccount(t *testing.T) {
	type fields struct {
		repo     pkg.IRepository
		gmailSvc *gmail.Service
		bundle   *i18n.Bundle
	}
	type args struct {
		ctx  context.Context
		data *model.IdentityAccount
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "general test",
			fields: fields{
				repo:     suite.repo,
				gmailSvc: suite.gmailSvc,
				bundle:   suite.bundle,
			},
			args: args{
				ctx: context.TODO(),
				data: &model.IdentityAccount{
					Name:              "test name",
					Email:             "test email",
					Phone:             "",
					PhoneAreaCode:     "",
					Password:          "",
					PhoneVerifyStatus: 1,
					EmailVerifyStatus: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test exists email",
			fields: fields{
				repo:     suite.repo,
				gmailSvc: suite.gmailSvc,
				bundle:   suite.bundle,
			},
			args: args{
				ctx: context.TODO(),
				data: &model.IdentityAccount{
					Name:              "test name",
					Email:             "test already exists email",
					Phone:             "",
					PhoneAreaCode:     "",
					Password:          "",
					PhoneVerifyStatus: 1,
					EmailVerifyStatus: 1,
				},
			},
			wantErr: true,
		},
		{
			name: "test exists phone",
			fields: fields{
				repo:     suite.repo,
				gmailSvc: suite.gmailSvc,
				bundle:   suite.bundle,
			},
			args: args{
				ctx: context.TODO(),
				data: &model.IdentityAccount{
					Name:              "test name",
					Email:             "",
					Phone:             "test already exists phone",
					PhoneAreaCode:     "test already exists phone area code",
					Password:          "",
					PhoneVerifyStatus: 1,
					EmailVerifyStatus: 1,
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo:     tt.fields.repo,
				gmailSvc: tt.fields.gmailSvc,
				bundle:   tt.fields.bundle,
			}

			if err := s.CreateIdentityAccount(tt.args.ctx, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("service.CreateIdentityAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
