package cauth

// Config configures the cauth module
type Config struct {
	VerificationCodeLen      uint   `toml:"verification_code_len" default:"6"`
	VerificationEmailSubject string `toml:"verification_email_subject" default:"Verification Code"`
	VerificationEmailFrom    string `toml:"verification_email_from" default:"webmaster@example.com"`
}
