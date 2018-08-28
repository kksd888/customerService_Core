// 健康检查

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/logic"
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
	context.JSON(http.StatusOK, InitResponse{
		BaseResponse:       BaseResponse{},
		Mine:               InitMine{},
		InitOnlineCustomer: []InitOnlineCustomer{},
	})
}

// API全局响应基础结构
type BaseResponse struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Token string `json:"token"`
}

type InitResponse struct {
	BaseResponse
	Mine               InitMine             `json:"mine"`
	InitOnlineCustomer []InitOnlineCustomer `json:"init_online_customer"`
}
type InitMine struct {
	Id         string `json:"id"`
	UserName   string `json:"user_name"`
	HeadImgUrl string `json:"head_img_url"`
	Status     string `json:"status"`
}
type InitOnlineCustomer struct {
	Id                 string              `json:"id"`
	CustomerNickName   string              `json:"customer_nick_name"`
	CustomerHeadImgUrl string              `json:"customer_head_img_url"`
	CustomerMessages   []logic.RoomMessage `json:"customer_messages"`
}
