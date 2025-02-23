package logging

import (
	"github.com/sirupsen/logrus"
	"os"
)

func NewLogging(level string, format string) *logrus.Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)

	switch level {
	case "trace":
		log.SetLevel(logrus.TraceLevel)
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	case "fatal":
		log.SetLevel(logrus.FatalLevel)
	case "panic":
		log.SetLevel(logrus.PanicLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	switch format {
	case "json":
		log.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		fallthrough
	default:
		log.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
			DisableQuote:  true,
		})
	}

	return log
}
