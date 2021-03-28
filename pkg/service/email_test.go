package service

import (
	pkg "amazing_talker/pkg"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"google.golang.org/api/gmail/v1"
)

func Test_service_getEmailBodyMessage(t *testing.T) {
	type fields struct {
		repo     pkg.IRepository
		gmailSvc *gmail.Service
		bundle   *i18n.Bundle
	}
	type args struct {
		lang string
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "test chinese",
			fields: fields{
				repo:     suite.repo,
				gmailSvc: suite.gmailSvc,
				bundle:   suite.bundle,
			},
			args: args{
				lang: "ch",
				name: "test name",
			},
			want: "test name 您好! 歡迎加入 AmazingTalker",
		},
		{
			name: "test english",
			fields: fields{
				repo:     suite.repo,
				gmailSvc: suite.gmailSvc,
				bundle:   suite.bundle,
			},
			args: args{
				lang: "en",
				name: "test name",
			},
			want: "Hi test name ! Welcome to AmazingTalker.",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &service{
				repo:     tt.fields.repo,
				gmailSvc: tt.fields.gmailSvc,
				bundle:   tt.fields.bundle,
			}
			if got := s.getEmailBodyMessage(tt.args.lang, tt.args.name); got != tt.want {
				t.Errorf("service.getEmailBodyMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
