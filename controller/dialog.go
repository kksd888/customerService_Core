// 对话管理

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
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
// @Router /v1/dialog/:kfId/list [get]
func (c *DialogController) List(context *gin.Context) {
	kfId := context.Param("kfId")
	if kfId == "" {
		ReturnErrInfo(context, errors.New("参数不能为空"))
	}
	iKfId, err := strconv.Atoi(kfId)
	if err != nil {
		ReturnErrInfo(context, errors.New("客服编号异常"))
	}

	customer := model.MessageLinkCustomer{Message: model.Message{KfId: iKfId}}
	messages, e := customer.WaitReply()
	ReturnErrInfo(context, e)

	context.JSON(http.StatusOK, messages)
}

// @Summary 客服接入用户，确认该用户所有未读消息
// @Description 客服接入用户，确认该用户所有未读消息
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/dialog/access [post]
func (c *DialogController) Access(context *gin.Context) {
	var aRequest AccessRequest
	context.Bind(&aRequest)
	roomKf, _ := handle.AuthToken2Model(context)

	// 所有对应客户的ACK消息确认
	for _, v := range aRequest.CustomerIds {
		model.Message{CustomerToken: v, KfId: roomKf.KfId}.AccessAck()
	}

	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		logic.KfAccess(aRequest.CustomerIds, logic.RoomKf{})
	}
}

type AccessRequest struct {
	CustomerIds []string `json:"customer_ids"`
}
