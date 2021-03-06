package cmailer

import (
	"context"

	"github.com/tusharsoni/copper/clogger"
)

type LogMailer struct {
	logger clogger.Logger
}

func NewLogMailer(logger clogger.Logger) Mailer {
	return &LogMailer{
		logger: logger,
	}
}

func (m *LogMailer) Send(ctx context.Context, p SendParams) error {
	m.logger.WithTags(map[string]interface{}{
		"from":      p.From,
		"to":        p.To,
		"subject":   p.Subject,
		"htmlBody":  p.HTMLBody,
		"plainBody": p.PlainBody,
	}).Info("Send plain email")

	return nil
}
