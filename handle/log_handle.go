package handle

import (
	"github.com/gin-gonic/gin"
)

func LogMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
