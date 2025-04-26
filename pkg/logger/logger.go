package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func InitLogger() {
	Log = logrus.New()
	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.DebugLevel) // Чтобы показывались debug/info/error
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
