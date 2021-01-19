package utils

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
)

// GenJwtToken 使用key作为签名秘钥来生成jwt token，并且在token里包含了自定义数据data
// 当配置文件中包含jwt.exp字段，并且该字段大于零时，在jwt token里会加入exp来表示token过期时间
// 单位是小时
func GenJwtToken(key string, data map[string]interface{}) (string, error) {
	mapClaims := jwt.MapClaims{}
	if exp := viper.GetInt("jwt.exp"); exp > 0 {
		mapClaims["exp"] = time.Now().Add(time.Hour * time.Duration(exp))
	}

	for k, v := range data {
		mapClaims[k] = v
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenStr, err := t.SignedString([]byte(key))
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// ParseJwtToken 解析jwt token，返回*jwt.Token对象
func ParseJwtToken(token string, key string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})

	if err != nil {
		switch err.(type) {
		case *jwt.ValidationError:
			vErr := err.(*jwt.ValidationError)
			switch vErr.Errors {
			case jwt.ValidationErrorExpired:
				return nil, fmt.Errorf("Token expired")
			default:
				return nil, fmt.Errorf("Parse token error")
			}
		default:
			return nil, fmt.Errorf("Invalid token")
		}
	}
	return t, nil
	//token2User[t.Raw] = t.Claims.(jwt.MapClaims)["phone"].(string)
}

// GetValueFromJwtToken 从ctx中解析jwt token(token存在于请求头的"jwt.headerName"中)
// 并从解析出来的对象中取得mapkey的值(interface{}类型) v，进一步可使用v.(someType)来转成具体类型
func GetValueFromJwtToken(ctx *gin.Context, mapkey string) (interface{}, error) {

	token := ctx.Request.Header.Get(viper.GetString("jwt.headerName"))
	t, err := ParseJwtToken(token, viper.GetString("jwt.mapkey"))
	if err != nil {
		msg := fmt.Sprintf("Parse jwt token error: %v", err)
		logrus.Error(msg)
		return nil, err
	}
	value := t.Claims.(jwt.MapClaims)[mapkey]
	return value, nil
}
