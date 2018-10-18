package handle

import (
	"encoding/base64"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// 开放API鉴权
func OpenApiOauthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err           error
			customerId    string
			authorization = c.Request.Header.Get("authorization")
		)

		// 检测授权先决条件
		if authorization == "" {
			c.JSON(http.StatusUnauthorized, common.BaseOutput{
				Code:   common.HeaderBad,
				ErrMsg: common.AUTHORIZATION_REQUIRED.Error(),
			})
			log.Error(common.AUTHORIZATION_REQUIRED)

			c.Abort()
			return
		}

		// 检测授权正确性
		if customerId, err = OpenAuthToken2Model(authorization); err != nil {
			c.JSON(http.StatusUnauthorized, common.BaseOutput{
				Code:   common.SignBad,
				ErrMsg: err.Error(),
			})
			log.Error(err)

			c.Abort()
			return
		} else {
			c.Set("CID", customerId)
		}

		c.Next()
	}
}

// authorization 解码
func OpenAuthToken2Model(authorization string) (customerId string, err error) {
	var (
		aes         = common.AesEncrypt{}
		decodeBytes []byte
	)

	// base64解析
	if decodeBytes, err = base64.StdEncoding.DecodeString(authorization); err != nil {
		log.WithFields(log.Fields{"default": err}).Error("base64解码错误")
		err = common.AUTHORIZATION_FAILED
		return
	}

	// aes解码
	if bytes, err := aes.Decrypt(decodeBytes); err != nil {
		log.WithFields(log.Fields{"default": err}).Warn("authorization aes解析错误")
		err = common.AUTHORIZATION_FAILED
	} else {
		customerId = string(bytes)
		if customerId == "" {
			log.WithFields(log.Fields{"default": err}).Warn("authorization aes解析错误")
			err = common.AUTHORIZATION_FAILED
		}
	}
	return
}

// authorization 编码
func OpenMake2Auth(customerId string) (string, error) {
	encrypt := common.AesEncrypt{}
	byteInfo, err := encrypt.Encrypt([]byte(customerId))
	if err != nil {
		log.WithFields(log.Fields{"default": err}).Warn("authorization 编码错误")
	}

	return base64.StdEncoding.EncodeToString(byteInfo), err
}
