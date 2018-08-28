package handle

import (
	"encoding/json"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"github.com/gin-gonic/gin"
)

// OAuth2.0 授权认证
func OauthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("authentication")
		if token == "" {
			c.JSON(401, "API token required")
			c.Abort()
			return
		}

		if err := AuthToken2Model(token, new(logic.RoomKf)); err != nil {
			c.JSON(401, err.Error())
			c.Abort()
			return
		}

		c.Next()
	}
}

// 鉴权Token解码为模型
func AuthToken2Model(token string, roomKf *logic.RoomKf) (err error) {
	aes := common.AesEncrypt{}
	if bytes, err := aes.Decrypt([]byte(token)); err != nil {
		err = errors.New("API token Authentication failed")
	} else {
		var roomKf logic.RoomKf
		if err := json.Unmarshal(bytes, &roomKf); err != nil {
			err = errors.New("API token Authentication failed")
		}
		if &roomKf == nil {
			err = errors.New("API token Authentication failed")
		}
	}
	return
}
