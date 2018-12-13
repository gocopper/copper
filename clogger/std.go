package clogger

import (
	"encoding/json"
	"log"
)

type StdLogger struct{}

func NewStdLogger() Logger {
	return &StdLogger{}
}

func (s *StdLogger) Debugw(msg string, tags map[string]string) {
	s.logw(DEBUG, msg, tags)
}

func (s *StdLogger) Infow(msg string, tags map[string]string) {
	s.logw(INFO, msg, tags)
}

func (s *StdLogger) Warnw(msg string, tags map[string]string) {
	s.logw(WARN, msg, tags)
}

func (s *StdLogger) Errorw(msg string, tags map[string]string) {
	s.logw(ERROR, msg, tags)
}

func (*StdLogger) logw(lvl Level, msg string, tags map[string]string) {
	if tags == nil {
		log.Printf("[%s] %s", string(lvl), msg)
		return
	}

	tagsJSON, _ := json.Marshal(tags)
	log.Printf("[%s] %s %s", string(lvl), msg, string(tagsJSON))
}
