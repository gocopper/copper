package emailotp

type Config struct {
	RequiresVerification bool
	VerificationEmail    EmailConfig
}

type EmailConfig struct {
	From         string
	Subject      string
	BodyTemplate string
}

func GetDefaultConfig() Config {
	return Config{
		RequiresVerification: false,
		VerificationEmail: EmailConfig{
			From:         "info@webmaster",
			Subject:      "Verify your account",
			BodyTemplate: "Your verification code is {{.VerificationCode}}",
		},
	}
}
