package utils

import (
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Fields Fields type
type Fields = logrus.Fields

func init() {
	log.Formatter = &logrus.JSONFormatter{}
}

// GetLog ...
func GetLog() *logrus.Logger {
	return log
}
