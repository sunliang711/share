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
		c.JSON(200, types.PushResp{Code: 1, Msg: msg})
		logrus.Errorf(msg)
		return
	}
	// check req.Content size
	pushLimit := viper.GetInt("business.push_limit")
	if len(req.Content) > pushLimit {
		msg := fmt.Sprintf("request content exceed %d bytes", pushLimit)
		c.JSON(200, types.PushResp{Code: 2, Msg: msg})
		return
	}

	logrus.Debugf("Request comming: %+v", req)
	md5Value := utils.Md5([]byte(req.Content))
	md5Str := hex.EncodeToString(md5Value)
	logrus.Debugf("Request hash: %v", md5Str)

	//hash prefix collision
	n := viper.GetInt("business.hash_min_len")
	newKey := ""
	for n < len(md5Str) {
		fields, err := database.Rdb.HMGet(database.Ctx, md5Str[:n], "type").Result()
		if err != nil {
			msg := "Internal redis error"
			c.JSON(200, types.PushResp{Code: 100, Msg: msg})
			logrus.Errorf(msg)
			return
		}
		// Not collision
		if fields[0] == nil {
			newKey = md5Str[:n]
			break
		}
		logrus.Debugf("hash collision: %v", md5Str[:n])
		n++
	}
	logrus.Debugf("Result key: %v", newKey)
	// md5 -> {type:xx, name: yy, content: zz}
	if err := database.Rdb.HSet(database.Ctx, newKey, "type", fmt.Sprintf("%d", types.TextType), "name", "", "content", string(req.Content)).Err(); err != nil {
		msg := "internal redis error"
		c.JSON(200, types.PushResp{Code: 100, Msg: msg})
		logrus.Errorf(msg)
		return
	}

	if err := database.Rdb.Expire(database.Ctx, newKey, time.Second*time.Duration(viper.GetInt64("business.ttl"))).Err(); err != nil {
		msg := "internal redis error"
		c.JSON(200, types.PushResp{Code: 101, Msg: msg})
		logrus.Errorf(msg)
		return
	}

	// c.JSON(200, types.Resp{Code: 0, Msg: "OK", Data: gin.H{"key": md5Str, "ttl": viper.GetInt64("business.ttl")}})
	c.JSON(200, types.PushResp{Code: 0, Msg: "OK", Key: newKey, TTL: viper.GetInt64("business.ttl")})
}

func Pull(c *gin.Context) {
	logrus.Infof("Content type: %v\n", c.ContentType())
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
