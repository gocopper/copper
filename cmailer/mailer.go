package cmailer

import "context"

type Mailer interface {
	SendPlain(ctx context.Context, from, to, subject, body string) (confirmation string, err error)
	SendHTML(ctx context.Context, from, to, subject, body string) (confirmation string, err error)
}
