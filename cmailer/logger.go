package cmailer

import (
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

func (m *LogMailer) SendPlain(from, to, subject, body string) (confirmation string, err error) {
	m.logger.WithTags(map[string]interface{}{
		"from":    from,
		"to":      to,
		"subject": subject,
		"body":    body,
	}).Info("Send plan email")

	return crandom.GenerateRandomString(6), nil
}

func (m *LogMailer) SendHTML(from, to, subject, body string) (confirmation string, err error) {
	m.logger.WithTags(map[string]interface{}{
		"from":    from,
		"to":      to,
		"subject": subject,
		"body":    body,
	}).Info("Send HTML email")

	return crandom.GenerateRandomString(6), nil
}
