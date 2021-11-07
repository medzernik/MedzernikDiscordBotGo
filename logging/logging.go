// Package logging serves to log all messages of the bot into the logs.log file
package logging

import (
	"github.com/medzernik/SlovakiaDiscordBotGo/config"
	"github.com/sirupsen/logrus"
	"os"
)

var Log = logrus.New()

func StartLogging() {

	Log.SetFormatter(&logrus.JSONFormatter{})

	//Set Log Level
	switch config.Cfg.ServerInfo.LogLevel {
	case "0":
		Log.SetLevel(logrus.TraceLevel)
	case "1":
		Log.SetLevel(logrus.DebugLevel)
	case "2":
		Log.SetLevel(logrus.InfoLevel)
	case "3":
		Log.SetLevel(logrus.WarnLevel)
	case "4":
		Log.SetLevel(logrus.ErrorLevel)
	case "5":
		Log.SetLevel(logrus.FatalLevel)
	default:
		Log.SetLevel(logrus.InfoLevel)

	}

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		Log.Out = file
	} else {
		Log.Info("Failed to log to file, using default stderr")
	}

	Log.Infof("Starting the bot...")

}
