package godid

import (
	"io/ioutil"
	"os"

	"github.com/mitchellh/go-homedir"

	"github.com/sirupsen/logrus"
)

const (
	logPath = workDir + "godid.log"
)

var (
	logger *logrus.Logger
)

func getLogger() *logrus.Logger {
	if logger != nil {
		return logger
	}
	logger = logrus.New()
	logger.SetOutput(ioutil.Discard)
	logger.Formatter = &logrus.TextFormatter{}
	if os.Getenv("GODID_TEST") != "" {
		return logger
	}
	if path, err := homedir.Expand(logPath); err == nil {
		if f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600); err == nil {
			logger.SetOutput(f)
			if os.Getenv("GODID_DEBUG") != "" {
				logger.SetLevel(logrus.DebugLevel)
			} else {
				logger.SetLevel(logrus.ErrorLevel)
			}
		}
	}
	return logger
}
