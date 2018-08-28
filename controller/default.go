// 健康检查

package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type DefaultController struct {
}

func InitHealth() *DefaultController {
	return &DefaultController{}
}

// @Summary 健康检查
// @Description 应用程序健康检查接口
// @Tags Default
// @Accept json
// @Produce json
// @Success 200 {string} json ""
// @Router /v1/health [get]
func (c *DefaultController) Health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
}

// @Summary 长连接监听
// @Description 长连接监听，支持和游览器通过WebSocket通信
// @Tags Default
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/ws [get]
func (c *DefaultController) Ws(context *gin.Context) {
}

// @Summary 系统初始化
// @Description 在线客服系统进行初始化
// @Tags Default
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/init [get]
func (c *DefaultController) Init(context *gin.Context) {
}
