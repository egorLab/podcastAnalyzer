package logging

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"time"
)

var Logger = logrus.New()

func CheckErr(err error, details string) {
	if err != nil {
		Logger.WithFields(logrus.Fields{
			"err":     err,
			"details": details,
		}).Error()
	}
}

func InitLogger(debug bool) {
	var logLevel = logrus.InfoLevel
	if debug {
		logLevel = logrus.DebugLevel
	}

	rotateFileHook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename:   "../logrus.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, //days
		Level:      logLevel,
		Formatter: &logrus.JSONFormatter{
			TimestampFormat: time.RFC822,
		},
	})

	if err != nil {
		Logger.Fatalf("Failed to initialize file rotate hook: %v", err)
	}

	Logger.SetLevel(logLevel)
	Logger.SetOutput(colorable.NewColorableStdout())
	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC822,
	})
	Logger.AddHook(rotateFileHook)
}

//func InitLogger() {
//	// create logger
//	file, err := os.OpenFile("../logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
//
//	if err == nil {
//		mw := io.MultiWriter(os.Stdout, file)
//		Logger.SetOutput(mw)
//	} else {
//		Logger.Info("Failed to log to file, using default stderr")
//	}
//
//	Logger.Info("Started logging...")
//}
