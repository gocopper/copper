package ctexter

import "github.com/tusharsoni/copper/clogger"

type loggerSvc struct {
	logger clogger.Logger
}

func newLoggerSvc(logger clogger.Logger) Svc {
	return &loggerSvc{logger: logger}
}

func (s *loggerSvc) SendSMS(phoneNumber, message string) (confirmationCode string, err error) {
	s.logger.WithTags(map[string]interface{}{
		"phoneNumber": phoneNumber,
		"message":     message,
	}).Info("texter::SendSMS")

	return "", nil
}
