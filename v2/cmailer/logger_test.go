package cmailer_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tusharsoni/copper/v2/clogger"
	"github.com/tusharsoni/copper/v2/cmailer"
)

func TestNewLogMailer(t *testing.T) {
	t.Parallel()

	mailer := cmailer.NewLogMailer(nil)

	_, ok := mailer.(cmailer.Mailer)

	assert.NotNil(t, mailer)
	assert.True(t, ok)
}

func TestLogMailer_Send(t *testing.T) {
	t.Parallel()

	var (
		logs   = make([]clogger.RecordedLog, 0)
		logger = clogger.NewRecorder(&logs)
		mailer = cmailer.NewLogMailer(logger)

		htmlBody  = "html-body"
		plainBody = "plain-body"
	)

	err := mailer.Send(context.Background(), cmailer.SendParams{
		From:      "from@test",
		To:        []string{"to@test"},
		Subject:   "test subject",
		HTMLBody:  &htmlBody,
		PlainBody: &plainBody,
	})

	assert.Nil(t, err)
	assert.Equal(t, clogger.LevelInfo, logs[0].Level)
	assert.Equal(t, "Send email", logs[0].Msg)
	assert.Equal(t, map[string]interface{}{
		"from":      "from@test",
		"to":        []string{"to@test"},
		"subject":   "test subject",
		"htmlBody":  &htmlBody,
		"plainBody": &plainBody,
	}, logs[0].Tags)
}
