package cmailer

import (
	"context"

	"github.com/tusharsoni/copper/clogger"
	"github.com/tusharsoni/copper/crandom"
)

type LogMailer struct {
	logger clogger.Logger
}

func NewLogMailer(logger clogger.Logger) Mailer {
	return &LogMailer{
		logger: logger,
	}
}

func (m *LogMailer) SendPlain(ctx context.Context, from, to, subject, body string) (confirmation string, err error) {
	m.logger.WithTags(map[string]interface{}{
		"from":    from,
		"to":      to,
		"subject": subject,
		"body":    body,
	}).Info("Send plain email")

	return crandom.GenerateRandomString(6), nil
}

func (m *LogMailer) SendHTML(ctx context.Context, from, to, subject, body string) (confirmation string, err error) {
	m.logger.WithTags(map[string]interface{}{
		"from":    from,
		"to":      to,
		"subject": subject,
		"body":    body,
	}).Info("Send HTML email")

	return crandom.GenerateRandomString(6), nil
}
