// 访客操作

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/kf"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type CustomerController struct {
	wxContext *wechat.Wechat
	rooms     map[string]*logic.Room
}

func InitCustomer(wxContext *wechat.Wechat, rooms map[string]*logic.Room) *CustomerController {
	return &CustomerController{wxContext, rooms}
}

// @Summary 获取一个用户的聊天记录
// @Description 获取一个用户的聊天记录
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Param id path int true "客户 ID"
// @Success 200 {string} json ""
// @Router /v1/customer/{customerId}/history [get]
func (c *CustomerController) History(context *gin.Context) {
}

// @Summary 客服发送消息给客户
// @Description 客服发送消息给客户
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Param id path int true "客户 ID"
// @Success 200 {string} json ""
// @Router /v1/customer/message [post]
func (c *CustomerController) SendMessage(context *gin.Context) {
	var sendRequest SendMessageRequest
	if bindErr := context.Bind(&sendRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	roomKf, err := handle.AuthToken2Model(context)
	ReturnErrInfo(context, err)

	model.Message{
		CustomerToken: sendRequest.CustomerId,
		KfId:          roomKf.KfId,
		MsgType:       sendRequest.MsgType,
		Msg:           sendRequest.Msg,
		OperCode:      common.MessageFromKf,
		KfAck:         true,
	}.Insert()

	msgResponse, err := c.wxContext.GetKf().Send(kf.KfSendMsgRequest{
		ToUser:  sendRequest.CustomerId,
		MsgType: sendRequest.MsgType,
		Text: message.Text{
			Content: sendRequest.Msg,
		},
	})
	ReturnErrInfo(context, err)

	if msgResponse.ErrCode == 0 {
		ReturnSuccessInfo(context)
	} else {
		ReturnErrInfo(context, errors.New("发送消息失败"))
	}
}

// 发送消息
type SendMessageRequest struct {
	CustomerId string `json:"customer_id"`
	MsgType    string `json:"msg_type"`
	Msg        string `json:"msg"`
}
