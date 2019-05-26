package cauth

import (
	"golang.org/x/crypto/bcrypt"
)

// Config for services in the cauth package that can be provided using Fx.
type Config struct {
	VerificationCodeLen   uint
	ResetPasswordTokenLen uint
	PasswordHashCost      int
	VerificationEmail     EmailConfig
	ResetPasswordEmail    EmailConfig
	SessionTokenLen       uint
	AdminEmail            string
}

// EmailConfig can be used to configure the email that is sent during various authentication flows such as user
// verification and reset password.
type EmailConfig struct {
	From         string
	Subject      string
	BodyTemplate string
}

// GetDefaultConfig provides a default set of config with sane defaults.
// Note: To send verification emails successfully, override the VerificationEmail.From property with an authorized
// email address.
func GetDefaultConfig() Config {
	return Config{
		VerificationCodeLen:   6,
		SessionTokenLen:       72,
		ResetPasswordTokenLen: 8,
		PasswordHashCost:      bcrypt.DefaultCost,
		VerificationEmail: EmailConfig{
			From:         "info@webmaster",
			Subject:      "Verify your account",
			BodyTemplate: "Your verification code is {{.VerificationCode}}",
		},
		ResetPasswordEmail: EmailConfig{
			From:         "info@webmaster",
			Subject:      "Reset password",
			BodyTemplate: "Your password has been reset to {{.ResetToken}}. Please login and change your password immediately.",
		},
		AdminEmail: "admin@webmaster",
	}
}
