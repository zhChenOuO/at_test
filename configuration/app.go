package configuration

// App ..
type App struct {
	JwtSecrets   string `json:"jwt_secrets" mapstructure:"jwt_secrets"`
	JwtExpireSec int    `json:"jwt_expire_sec" mapstructure:"jwt_expire_sec"`

	Scheduler struct {
		VerifyIdentityAccountEmailFreq string `mapstructure:"verify_identity_account_email_freq"`
		VerifyIdentityAccountPhoneFreq string `mapstructure:"verify_identity_account_phone_freq"`
	}
}
