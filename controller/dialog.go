// 对话管理

package controller

import (
	"encoding/json"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/kf"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type DialogController struct {
	wxContext *wechat.Wechat
	rooms     map[string]*logic.Room
}

func InitDialog(wxContext *wechat.Wechat, rooms map[string]*logic.Room) *DialogController {
	return &DialogController{wxContext: wxContext, rooms: rooms}
}

// @Summary 待接入列表
// @Description 待接入列表
// @Tags WaitQueue
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/wait_queue [get]
func (c *DialogController) Queue(context *gin.Context) {
	if waitQueueRooms, err := logic.GetWaitQueue(); err != nil {
		ReturnErrInfo(context, err)
	} else {
		var waitQueues []WaitQueueResponse
		for _, value := range waitQueueRooms {
			waitQueues = append(waitQueues, WaitQueueResponse{
				CustomerId:         value.CustomerId,
				CustomerNickName:   value.CustomerNickName,
				CustomerHeadImgUrl: value.CustomerHeadImgUrl,
				//Messages:           value.CustomerMsgs,
				PreviousKf: WaitQueuePreviousKf{},
			})
		}
		context.JSON(http.StatusOK, waitQueues)
	}
}

// @Summary 会话确认应答
// @Description 会话确认应答
// @Tags WaitQueue
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /v1/wait_queue/access [post]
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

// @Summary 获取待回复消息列表
// @Description 获取待回复消息列表
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/dialog [get]
func (c *DialogController) List(context *gin.Context) {
	roomKf, _ := handle.AuthToken2Model(context)

	customer := model.MessageLinkCustomer{Message: model.Message{KfId: roomKf.KfId}}
	messages, e := customer.WaitReply()
	ReturnErrInfo(context, e)

	bytes, _ := json.Marshal(messages)
	log.Println(string(bytes))

	context.JSON(http.StatusOK, messages)
}

// @Summary 确认已读
// @Description 确认已读
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
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

// @Summary 获取聊天记录
// @Description 获取聊天记录
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Param customerId path int true "客户 ID"
// @Success 200 {string} json ""
// @Router /v1/dialog/{customerId} [get]
func (c *DialogController) History(context *gin.Context) {
}

// @Summary 发送消息
// @Description 发送消息
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /v1/dialog [post]
func (c *DialogController) SendMessage(context *gin.Context) {
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

// 访客队列响应
type WaitQueueResponse struct {
	CustomerId         string
	CustomerNickName   string
	CustomerHeadImgUrl string
	Messages           []*logic.RoomMessage
	PreviousKf         WaitQueuePreviousKf
}
type WaitQueuePreviousKf struct {
	KfId     string
	KfName   string
	LastTime time.Time
}

type CustomerIdsRequest struct {
	CustomerIds []string `json:"customer_ids"`
}

// 发送消息
type SendMessageRequest struct {
	CustomerId string `json:"customer_id"`
	MsgType    string `json:"msg_type"`
	Msg        string `json:"msg"`
}
