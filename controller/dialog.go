// 对话管理

package controller

import (
	"encoding/json"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/model"
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
// @Router /v1/dialog/list [get]
func (c *DialogController) List(context *gin.Context) {
	roomKf, _ := handle.AuthToken2Model(context)

	customer := model.MessageLinkCustomer{Message: model.Message{KfId: roomKf.KfId}}
	messages, e := customer.WaitReply()
	ReturnErrInfo(context, e)

	bytes, _ := json.Marshal(messages)
	log.Println(string(bytes))

	context.JSON(http.StatusOK, messages)
}

// @Summary 客服接入用户，客服加入会话房间
// @Description 客服接入用户，客服加入会话房间
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/dialog/access [post]
func (c *DialogController) Access(context *gin.Context) {
	var aRequest CustomerIdsRequest
	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	roomKf, _ := handle.AuthToken2Model(context)

	for _, v := range aRequest.CustomerIds {
		// 客服加入聊天房间
		room, _ := logic.InitRoom(v)
		room.RoomKf = logic.RoomKf{
			KfId:         roomKf.KfId,
			KfName:       roomKf.KfName,
			KfHeadImgUrl: roomKf.KfHeadImgUrl,
			KfStatus:     common.KF_ONLINE,
		}

		// 更新所有指定客户的KfId
		model.Message{CustomerToken: v, KfId: roomKf.KfId}.Access()
	}

	ReturnSuccessInfo(context)
}

// @Summary 确认已读用户的消息
// @Description 确认已读用户的消息
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/dialog/ack [put]
func (c *DialogController) Ack(context *gin.Context) {
	var aRequest CustomerIdsRequest
	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}
	roomKf, _ := handle.AuthToken2Model(context)

	for _, v := range aRequest.CustomerIds {
		model.Message{CustomerToken: v, KfId: roomKf.KfId, KfAck: true}.Ack()
	}

	ReturnSuccessInfo(context)
}

type CustomerIdsRequest struct {
	CustomerIds []string `json:"customer_ids"`
}
