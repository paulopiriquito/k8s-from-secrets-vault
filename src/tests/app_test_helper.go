package tests

import (
	"github.com/sirupsen/logrus"
	"os"
	"testing"
)

func setupLogger(t *testing.T) *logrus.Logger {
	t.Helper()

	log := logrus.New()
	log.Out = os.Stdout
	log.Formatter = &logrus.JSONFormatter{}
	return log
}
