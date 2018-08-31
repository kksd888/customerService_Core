package handle

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// OAuth2.0 授权认证
func OauthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authentication")
		if token == "" {
			c.JSON(http.StatusUnauthorized, "API token required")
			c.Abort()
			return
		}

		if roomKf, err := AuthToken2Model(c); err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		} else {
			c.Set("roomKf", roomKf)
		}
		c.Next()
	}
}

// 鉴权Token解码为模型
func AuthToken2Model(c *gin.Context) (roomKf *model.RoomKf, err error) {
	token := c.Request.Header.Get("authentication")
	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	aes := common.AesEncrypt{}
	if bytes, err := aes.Decrypt(decodeBytes); err != nil {
		err = errors.New("API token Authentication failed")
	} else {
		if err := json.Unmarshal(bytes, &roomKf); err != nil {
			log.Printf("string json :%s, err %#v", string(bytes), err.Error())
			err = errors.New("API token Authentication failed")
		}
		if &roomKf == nil {
			err = errors.New("API token Authentication failed")
		}
	}
	return
}
