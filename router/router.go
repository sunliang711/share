package router

import (
	"time"

	"share/handlers"
	"share/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// StartServer starts gin server
func StartServer(addr string, tls bool, certFile string, keyFile string) {
	//MUST SetMode first
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(gin.Logger())

	corsCfg := cors.Config{
		AllowOrigins: viper.GetStringSlice("cors.allowOrigins"),
		AllowMethods: viper.GetStringSlice("cors.allowMethods"),
		AllowHeaders: viper.GetStringSlice("cors.allowHeaders"),
		MaxAge:       time.Second * time.Duration(viper.GetInt("cors.maxAge")),
	}
	logrus.Infof(utils.CorsConfigStringify(&corsCfg))

	router.Use(cors.New(corsCfg))

	// Put normal handlers below
	router.GET("/api/v1/health", handlers.Health)
	router.POST("/api/v1/push_text", handlers.PushText)
	router.POST("/api/v1/push_file", handlers.PushFile)
	router.GET("/api/v1/pull/:key", handlers.Pull)
	// router.GET("/api/v1/PATH", handlers.XXX)

	// Put need-auth handlers below
	// router.GET("/api/v1/PATH", middleware.Auth)
	// router.POST("/api/v1/PATH", middleware.Auth)

	logrus.Infof("Start server on %v, tls enabled: %v", addr, tls)
	if tls {
		logrus.Fatalln(router.RunTLS(addr, certFile, keyFile))
	} else {
		logrus.Fatalln(router.Run(addr))
	}

}
