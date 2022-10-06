package services

import "github.com/sirupsen/logrus"

func GetLogger() *logrus.Logger {
	var log = logrus.New()

	log.SetFormatter(&logrus.TextFormatter{
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
		PadLevelText:    true,
	})

	return log
}
