// 客服相关

package controller

import "github.com/gin-gonic/gin"

type ServerController struct {
}

func InitServer() *ServerController {
	return &ServerController{}
}

// @Summary 获取客服信息
// @Description 获取客服信息
// @Tags Server
// @Accept  json
// @Produce  json
// @Param id path int true "客服的ID"
// @Success 200 {string} json ""
// @Router /v1/server/{id} [get]
func (c *ServerController) Get(context *gin.Context) {
}

// @Summary 客服修改在线状态
// @Description 客服修改在线状态
// @Tags Server
// @Accept  json
// @Produce  json
// @Param id path int true "客服的ID"
// @Success 200 {string} json ""
// @Router /v1/server/{id}/status [put]
func (c *ServerController) ChangeStatus(context *gin.Context) {
}
