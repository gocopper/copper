package cauth

import "golang.org/x/crypto/bcrypt"

// Config for services in the cauth package that can be provided using Fx.
type Config struct {
	VerificationCodeLen uint
	PasswordHashCost    int
	VerificationEmail   VerificationEmailConfig
	SessionTokenLen     uint
}

// VerificationEmailConfig can be used to configure the email that is sent after a user signs up.
// It can be provided as part of the Config.
type VerificationEmailConfig struct {
	From         string
	Subject      string
	BodyTemplate string
}

// GetDefaultConfig provides a default set of config with sane defaults.
// Note: To send verification emails successfully, override the VerificationEmail.From property with an authorized
// email address.
func GetDefaultConfig() Config {
	return Config{
		VerificationCodeLen: 6,
		SessionTokenLen:     72,
		PasswordHashCost:    bcrypt.DefaultCost,
		VerificationEmail: VerificationEmailConfig{
			From:         "info@webmaster",
			Subject:      "Verify your account",
			BodyTemplate: "Your verification code is {{.VerificationCode}}",
		},
	}
}
