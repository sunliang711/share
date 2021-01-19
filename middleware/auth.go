package middleware

import (
	"share/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	tokenExpired = "Token expired"
)

func Auth(c *gin.Context) {
	token := c.Request.Header.Get(viper.GetString("jwt.headerName"))
	_, err := utils.ParseJwtToken(token, viper.GetString("jwt.key"))
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		logrus.Errorf("Parse jwt token error: %v", err.Error())
		c.Abort()
		return
	}
	c.Next()
}
