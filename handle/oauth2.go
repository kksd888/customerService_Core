package handle

import (
	"encoding/base64"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/gin-gonic/gin"
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

		if kfId, err := AuthToken2Model(c); err != nil {
			c.JSON(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		} else {
			c.Set("KFID", kfId)
		}
		c.Next()
	}
}

// 鉴权Token解码为模型
func AuthToken2Model(c *gin.Context) (kfId string, err error) {
	var (
		token = c.Request.Header.Get("authentication")
		aes   = common.AesEncrypt{}
	)

	decodeBytes, err := base64.StdEncoding.DecodeString(token)
	if bytes, err := aes.Decrypt(decodeBytes); err != nil {
		err = errors.New("API token Authentication failed")
	} else {
		kfId = string(bytes)
		if kfId == "" {
			err = errors.New("API token Authentication failed")
		}
	}
	return
}
