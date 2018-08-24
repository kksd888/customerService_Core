// 对话管理

package controller

import "github.com/gin-gonic/gin"

type DialogController struct {
}

func InitDialog() *DialogController {
	return &DialogController{}
}

func (c *DialogController) DialogInit(context *gin.Context) {
}

func (c *DialogController) List(context *gin.Context) {
}

func (c *DialogController) History(context *gin.Context) {
}

func (c *DialogController) Create(context *gin.Context) {
}

func (c *DialogController) SendMessage(context *gin.Context) {
}

func (c *DialogController) RecallMessage(context *gin.Context) {
}
