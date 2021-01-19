package config

import (
	"io"
	"os"

	"share/utils"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	logrus.Infof("config init()...")
	configPath := pflag.StringP("config", "c", "config.toml", "config file")
	pflag.Parse()
	logrus.Info("Config file: ", *configPath)

	viper.SetConfigFile(*configPath)
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatalf("Read config file: %v error: %v", *configPath, err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: viper.GetBool("log.showFullTime")})
	if viper.GetBool("log.reportCaller") {
		logrus.Info("Logrus: enable report caller")
		logrus.SetReportCaller(true)
	}
	loglevel := viper.GetString("log.level")
	logrus.Infoln("Log level: ", loglevel)
	logrus.SetLevel(utils.LogLevel(loglevel))

	var output io.Writer
	logfilePath := viper.GetString("log.logfile")
	if logfilePath != "" {
		handler, err := os.OpenFile(logfilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			logrus.Fatalf("Open logfile: %v error: %v", logfilePath, err)
		}
		logrus.Infof("Logfile path: %v", logfilePath)
		output = handler
	} else {
		output = os.Stderr
	}
	logrus.SetOutput(output)
}
