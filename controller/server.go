// 客服相关

package controller

import "github.com/gin-gonic/gin"

type ServerController struct {
}

func InitServer() *ServerController {
	return &ServerController{}
}

func (c *ServerController) Get(context *gin.Context) {
}

func (c *ServerController) ChangeStatus(context *gin.Context) {
}
