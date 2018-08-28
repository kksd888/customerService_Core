// 对话管理

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type DialogController struct {
	wxContext *wechat.Wechat
	rooms     map[string]*logic.Room
}

func InitDialog(wxContext *wechat.Wechat, rooms map[string]*logic.Room) *DialogController {
	return &DialogController{wxContext, rooms}
}

// @Summary 获取待回复消息列表
// @Description 获取待回复消息列表
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/dialog/:id/list [get]
func (c *DialogController) List(context *gin.Context) {
	dialogId := context.Param("dialogId")
	if dialogId == "" {
		context.JSON(http.StatusOK, gin.H{"msg": "参数不能为空"})
	}

	r, ok := c.rooms[dialogId]
	if !ok {
		log.Println("查询空异常")
	} else {
		log.Printf("RoomMaps内容：%#v", r)
	}

	context.JSON(http.StatusOK, r.CustomerMsgs)
}

// @Summary 客服接入用户
// @Description 客服接入用户，加入会话房间
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/dialog/access [post]
func (c *DialogController) Access(context *gin.Context) {
	var aRequest AccessRequest
	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		// todo 客服接入通话房间
		logic.KfAccess(aRequest.CustomerIds, logic.RoomKf{})
	}
}

type AccessRequest struct {
	CustomerIds []string `json:"customer_ids"`
}
