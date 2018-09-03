// 对话管理

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

type DialogController struct {
	db        *model.MongoDb
	wxContext *wechat.Wechat
}

func InitDialog(wxContext *wechat.Wechat, _db *model.MongoDb) *DialogController {
	return &DialogController{wxContext: wxContext, db: _db}
}

// @Summary 待接入列表
// @Description 待接入列表
// @Tags WaitQueue
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/wait_queue [get]
func (c *DialogController) Queue(context *gin.Context) {
	var (
		waitCustomer = []WaitCustomer{}
	)

	roomCollection := c.db.C("room")
	roomCollection.Find(bson.M{"roomkf.kfid": ""}).All(&waitCustomer)

	context.JSON(http.StatusOK, waitCustomer)
}

// @Summary 会话确认应答
// @Description 会话确认应答
// @Tags WaitQueue
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /v1/wait_queue/access [post]
func (c *DialogController) Access(context *gin.Context) {
	var (
		aRequest CustomerIdsRequest
		kf       model.Kf
	)

	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	kfId, _ := context.Get("KFID")

	roomCollection := c.db.C("room")

	kfCollection := c.db.C("kf")
	kfCollection.Find(bson.M{"id": kfId}).One(&kf)

	for _, v := range aRequest.CustomerIds {
		// 客服加入聊天房间
		roomKf := model.RoomKf{
			KfId:         kf.Id,
			KfName:       kf.NickName,
			KfHeadImgUrl: kf.HeadImgUrl,
			KfStatus:     common.KF_ONLINE,
		}
		roomCollection.Update(bson.M{"roomcustomer.customerid": v}, bson.M{"$set": bson.M{"roomkf": roomKf}})
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
	//roomKf, _ := handle.AuthToken2Model(context)
	//
	//customer := model.MessageLinkCustomer{Message: model.Message{KfId: roomKf.KfId}}
	//messages, e := customer.WaitReply()
	//ReturnErrInfo(context, e)
	//
	//bytes, _ := json.Marshal(messages)
	//log.Println(string(bytes))
	//
	//context.JSON(http.StatusOK, messages)
}

// @Summary 确认已读
// @Description 确认已读
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /v1/dialog/ack [put]
func (c *DialogController) Ack(context *gin.Context) {
	//var aRequest CustomerIdsRequest
	//if bindErr := context.BindJSON(&aRequest); bindErr != nil {
	//	ReturnErrInfo(context, bindErr)
	//}
	//roomKf, _ := handle.AuthToken2Model(context)
	//
	//for _, v := range aRequest.CustomerIds {
	//	model.Message{CustomerToken: v, KfId: roomKf.KfId, KfAck: true}.Ack()
	//}
	//
	//ReturnSuccessInfo(context)
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
	//var sendRequest SendMessageRequest
	//if bindErr := context.Bind(&sendRequest); bindErr != nil {
	//	ReturnErrInfo(context, bindErr)
	//}
	//
	//roomKf, err := handle.AuthToken2Model(context)
	//ReturnErrInfo(context, err)
	//
	//model.Message{
	//	CustomerToken: sendRequest.CustomerId,
	//	KfId:          roomKf.KfId,
	//	MsgType:       sendRequest.MsgType,
	//	Msg:           sendRequest.Msg,
	//	OperCode:      common.MessageFromKf,
	//	KfAck:         true,
	//}.Insert()
	//
	//msgResponse, err := c.wxContext.GetKf().Send(kf.KfSendMsgRequest{
	//	ToUser:  sendRequest.CustomerId,
	//	MsgType: sendRequest.MsgType,
	//	Text: message.Text{
	//		Content: sendRequest.Msg,
	//	},
	//})
	//ReturnErrInfo(context, err)
	//
	//if msgResponse.ErrCode == 0 {
	//	ReturnSuccessInfo(context)
	//} else {
	//	ReturnErrInfo(context, errors.New("发送消息失败"))
	//}
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
