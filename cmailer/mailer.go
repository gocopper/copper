package cmailer

import "context"

type SendParams struct {
	From    string
	To      string
	Subject string

	HTMLBody  *string
	PlainBody *string
}

type Mailer interface {
	Send(ctx context.Context, p SendParams) error
}
