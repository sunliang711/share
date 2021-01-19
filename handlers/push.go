package handlers

import (
	"encoding/hex"
	"fmt"
	"share/database"
	"share/types"
	"share/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func PushFile(c *gin.Context) {
}
func PushText(c *gin.Context) {
	// Only use Share.Content field
	var req types.Share
	if err := c.BindJSON(&req); err != nil {
		msg := "invalid request format"
		c.JSON(200, gin.H{"code": 1, "msg": msg})
		logrus.Errorf(msg)
		return
	}
	//TODO check req.Content size

	logrus.Debugf("req: %+v", req)
	md5Value := utils.Md5([]byte(req.Content))
	md5Str := hex.EncodeToString(md5Value)

	md5Str = md5Str[:viper.GetInt("business.md5PreLen")]

	// md5 -> {type:xx, name: yy, content: zz}
	if err := database.Rdb.HSet(database.Ctx, md5Str, "type", fmt.Sprintf("%d", types.TextType), "name", "", "content", string(req.Content)).Err(); err != nil {
		msg := "internal redis error"
		c.JSON(200, gin.H{"code": 100, "msg": msg})
		logrus.Errorf(msg)
		return
	}

	if err := database.Rdb.Expire(database.Ctx, md5Str, time.Second*time.Duration(viper.GetInt64("business.ttl"))).Err(); err != nil {
		msg := "internal redis error"
		c.JSON(200, gin.H{"code": 101, "msg": msg})
		logrus.Errorf(msg)
		return
	}

	c.JSON(200, types.Resp{Code: 0, Msg: "OK", Data: gin.H{"key": md5Str}})
}

func Pull(c *gin.Context) {
	key := c.Param("key")
	keys := []string{"type", "name", "content"}
	fields, err := database.Rdb.HMGet(database.Ctx, key, keys...).Result()
	if err != nil {
		msg := "internal redis error"
		logrus.Errorf(msg)
		c.JSON(200, types.Share{
			Code: 1,
			Msg:  msg,
			Type: types.InvalidType,
		})
		return
	}

	empty := true
	for _, field := range fields {
		if field != nil {
			empty = false
		}
	}
	if empty {
		msg := "Not found"
		logrus.Infof(msg)
		c.JSON(200, types.Share{
			Code: 1,
			Msg:  msg,
			Type: types.InvalidType,
		})
		return
	}

	typeString := fields[0].(string)
	typ, err := strconv.Atoi(typeString)
	if err != nil {
		msg := "Internal type not int??"
		logrus.Errorf(msg)
		c.JSON(200, types.Share{
			Code: 1,
			Msg:  msg,
			Type: types.InvalidType,
		})
	}

	switch types.ShareType(typ) {
	case types.TextType:
		content := fields[2].(string)
		c.JSON(200, types.Share{
			Code:    types.OK,
			Msg:     "OK",
			Type:    types.TextType,
			Content: content,
		})
	case types.FileType:
		c.JSON(200, types.Resp{
			Code: types.OK,
			Msg:  "TODO",
		})
	default:
		c.JSON(200, types.Resp{
			Code: 1,
			Msg:  "Cannot reach here",
		})
	}

}
