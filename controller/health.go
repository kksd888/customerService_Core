// 健康检查

package controller

import (
	"github.com/gin-gonic/gin"
	"time"
)

type HealthController struct {
}

func InitHealth() *HealthController {
	return &HealthController{}
}

// @Summary 健康检查
// @Description 应用程序健康检查接口
// @Accept json
// @Produce json
// @Success 200 {string} json ""
// @Router /health [get]
func (c *HealthController) Health(context *gin.Context) {
	context.JSON(200, gin.H{"status": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
}

// @Summary 长连接监听
// @Description 长连接监听，支持和游览器通过WebSocket通信
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /ws [get]
func (c *HealthController) Ws(context *gin.Context) {
}
