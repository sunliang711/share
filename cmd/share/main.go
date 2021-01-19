package main

import (
	"fmt"

	"share/database"
	"share/router"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	buildstamp string
	githash    string
)

func main() {
	logrus.Infof("main()")
	defer func() {
		database.Release()
	}()

	addr := fmt.Sprintf(":%d", viper.GetInt("server.port"))
	tls := viper.GetBool("tls.enable")
	certFile := viper.GetString("tls.certFile")
	keyFile := viper.GetString("tls.keyFile")
	router.StartServer(addr, tls, certFile, keyFile)
}
