// 对话管理

package controller

import (
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/kf"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
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
		waitCustomer   = []WaitCustomer{}
		roomCollection = c.db.C("room")
	)

	query := []bson.M{
		{
			"$match": bson.M{"room_kf.kf_id": ""},
		},
		{
			"$project": bson.M{
				"room_customer": 1,
				"room_messages": bson.M{"$slice": []interface{}{"$room_messages", 0, 5}},
			},
		},
	}
	roomCollection.Pipe(query).All(&waitCustomer)
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
		aRequest       CustomerIdsRequest
		kfModel        model.Kf
		kfId, _        = context.Get("KFID")
		roomCollection = c.db.C("room")
		kfCollection   = c.db.C("kf")
	)

	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	kfCollection.Find(bson.M{"id": kfId}).One(&kfModel)
	for _, v := range aRequest.CustomerIds {
		// 客服加入聊天房间
		roomKf := model.RoomKf{
			KfId:         kfModel.Id,
			KfName:       kfModel.NickName,
			KfHeadImgUrl: kfModel.HeadImgUrl,
			KfStatus:     common.KF_ONLINE,
		}
		if err := roomCollection.Update(bson.M{"room_customer.customer_id": v}, bson.M{"$set": bson.M{"room_kf": roomKf}}); err != nil {
			ReturnErrInfo(context, err)
		}
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
	var (
		waitCustomer   = []WaitCustomer{}
		kfId, _        = context.Get("KFID")
		roomCollection = c.db.C("room")
	)

	roomCollection.Find(bson.M{"room_kf.kf_id": kfId, "room_messages.ack": false}).All(&waitCustomer)

	context.JSON(http.StatusOK, waitCustomer)
}

// @Summary 确认已读
// @Description 确认已读
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /v1/dialog/ack [put]
func (c *DialogController) Ack(context *gin.Context) {
	var (
		aRequest       CustomerIdsRequest
		kfId, _        = context.Get("KFID")
		roomCollection = c.db.C("room")
	)
	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	for _, v := range aRequest.CustomerIds {
		if updateErr := roomCollection.Update(bson.M{"room_kf.kf_id": kfId, "room_customer.customer_id": v}, bson.M{"$set": bson.M{"room_messages.$[].ack": true}}); updateErr != nil {
			ReturnErrInfo(context, updateErr)
		}
	}
	ReturnSuccessInfo(context)
}

// @Summary 获取聊天记录
// @Description 获取聊天记录
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Param customerId path int true "客户 ID"
// @page customerId path int true "第几页"
// @limit customerId path int true "页容量"
// @Success 200 {string} json ""
// @Router /v1/dialog/{customerId}/{page}/{limit} [get]
func (c *DialogController) History(context *gin.Context) {
	var (
		roomHistory    RoomHistory
		customerId     = context.Param("customerId")
		strPage        = context.Param("page")
		strLimit       = context.Param("limit")
		roomCollection = c.db.C("room")
	)
	if customerId == "" {
		ReturnErrInfo(context, errors.New("缺少customerId"))
	}

	page, err := strconv.Atoi(strPage)
	if err != nil {
		ReturnErrInfo(context, errors.New("缺少page"))
	}
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		ReturnErrInfo(context, errors.New("缺少limit"))
	}

	query := []bson.M{
		{
			"$match": bson.M{"room_customer.customer_id": customerId},
		},
		{
			"$project": bson.M{
				"room_messages": bson.M{"$slice": []interface{}{"$room_messages", (page - 1) * limit, limit}},
			},
		},
	}
	if err := roomCollection.Pipe(query).One(&roomHistory); err != nil {
		ReturnErrInfo(context, err)
	}

	context.JSON(http.StatusOK, roomHistory)
}

// @Summary 发送消息
// @Description 发送消息
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /v1/dialog [post]
func (c *DialogController) SendMessage(context *gin.Context) {
	var (
		sendRequest    SendMessageRequest
		kfId, _        = context.Get("KFID")
		roomCollection = c.db.C("room")
	)
	if bindErr := context.Bind(&sendRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	roomCollection.Update(bson.M{"room_kf.kf_id": kfId, "room_customer.customer_id": sendRequest.CustomerId},
		bson.M{"$push": bson.M{"room_messages": &model.RoomMessage{
			Id:         common.GetNewUUID(),
			Type:       sendRequest.MsgType,
			Msg:        sendRequest.Msg,
			OperCode:   common.MessageFromKf,
			Ack:        true,
			CreateTime: time.Now(),
		}}})

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

type CustomerIdsRequest struct {
	CustomerIds []string `json:"customer_ids"`
}

type SendMessageRequest struct {
	CustomerId string `json:"customer_id"`
	MsgType    string `json:"msg_type"`
	Msg        string `json:"msg"`
}

type RoomHistory struct {
	RoomMessages []model.RoomMessage `bson:"room_messages" json:"room_messages"`
}
