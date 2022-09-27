package logger

import (
	"context"
	"os"
	"time"

	filename "github.com/keepeye/logrus-filename"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"github.com/trucktrace/internal/models"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

var MyUserContext context.Context

type CtxKey struct{}

var (
	infoLogFile  = "logs/info_" + time.Now().Format("02-01-2006") + ".log"
	errorLogFile = "logs/erro_" + time.Now().Format("02-01-2006") + ".log"
)

func InitLogRus() error {
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		err = os.Mkdir("logs", 0755)
		if err != nil {
			return err
		}
	}

	Formatter := new(logrus.TextFormatter)

	Formatter.TimestampFormat = "02-01-2006 15:04:05"
	Formatter.FullTimestamp = true
	logrus.SetFormatter(Formatter)

	filenameHook := filename.NewHook()
	filenameHook.Field = "line"

	logrus.AddHook(filenameHook)

	writerHook := lfshook.NewHook(
		lfshook.PathMap{
			logrus.ErrorLevel: errorLogFile,
			logrus.InfoLevel:  infoLogFile,
		},
		&prefixed.TextFormatter{
			FullTimestamp:    true,
			TimestampFormat:  "02-01-2006 15:04:05",
			ForceColors:      true,
			ForceFormatting:  true,
			DisableUppercase: false,
			DisableColors:    true,
		},
	)

	logrus.AddHook(writerHook)

	return nil
}

func InfoLogger(method string) *logrus.Entry {

	if MyUserContext == nil {
		return SystemLoggerInfo(method)
	} else {
		user := MyUserContext.Value(CtxKey{}).(models.User)

		return logrus.WithFields(logrus.Fields{
			"method":   method,
			"username": user.Username,
			"email":    user.Email,
		})
	}
}

func ErrorLogger(method, error string) *logrus.Entry {
	if MyUserContext == nil {
		return SystemLoggerError(method, error)
	} else {
		user := MyUserContext.Value(CtxKey{}).(models.User)
		return logrus.WithFields(logrus.Fields{
			"method":       method,
			"username":     user.Username,
			"email":        user.Email,
			"errorMessage": error,
		})
	}
}

func SystemLoggerError(method, error string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"method":       method,
		"errorMessage": error,
	})
}

func SystemLoggerInfo(method string) *logrus.Entry {
	return logrus.WithFields(logrus.Fields{
		"method": method,
	})
}
