package model

import (
	"amazing_talker/configuration"
	"amazing_talker/pkg/model/ctype"
	"testing"

	"gotest.tools/assert"
)

func TestIdentityAccount_CreateToken(t *testing.T) {
	type fields struct {
		ID            int
		Name          string
		Email         string
		Phone         string
		PhoneAreaCode string
	}
	type args struct {
		cfg *configuration.App
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "general test",
			fields: fields{
				ID:    1,
				Name:  "test name",
				Email: "test email",
				Phone: "test phone",
			},
			args: args{
				cfg: &configuration.App{
					JwtSecrets:   "abc",
					JwtExpireSec: 3600000,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := &IdentityAccount{
				ID:    tt.fields.ID,
				Name:  tt.fields.Name,
				Email: tt.fields.Email,
				Phone: tt.fields.Phone,
			}
			_, err := identity.CreateToken(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("IdentityAccount.CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestIdentityAccount_Init(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "general test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			identity := IdentityAccount{}
			identity.Init()

			assert.Equal(t, identity.EmailVerifyStatus, ctype.VerifyInit)
			assert.Equal(t, identity.PhoneVerifyStatus, ctype.VerifyInit)
		})
	}
}
