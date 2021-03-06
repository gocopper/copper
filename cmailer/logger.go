package cmailer

import (
	"context"

	"github.com/tusharsoni/copper/clogger"
)

// NewLogMailer creates an implementation of Mailer that logs
// all sends with the provided logger. Useful during dev as
// it requires no configuration.
func NewLogMailer(logger clogger.Logger) Mailer {
	return &logMailer{logger: logger}
}

type logMailer struct {
	logger clogger.Logger
}

func (m *logMailer) Send(ctx context.Context, p SendParams) error {
	m.logger.WithTags(map[string]interface{}{
		"from":      p.From,
		"to":        p.To,
		"subject":   p.Subject,
		"htmlBody":  p.HTMLBody,
		"plainBody": p.PlainBody,
	}).Info("Send email")

	return nil
}
