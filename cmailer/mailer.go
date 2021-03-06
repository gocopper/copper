package cmailer

import "context"

// Mailer provides methods to send emails.
type Mailer interface {
	Send(ctx context.Context, p SendParams) error
}

// SendParams holds data needed to send an email using Mailer.
type SendParams struct {
	From    string
	To      []string
	Subject string

	HTMLBody  *string
	PlainBody *string
}
