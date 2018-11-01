// 对话管理

package admin

import (
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/kf"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type DialogController struct {
	wxContext *wechat.Wechat
}

func NewDialog(wxContext *wechat.Wechat) *DialogController {
	return &DialogController{wxContext: wxContext}
}

// @Summary 获取待回复消息列表 (5s轮询一次)
// @Description 获取待回复消息列表 (5s轮询一次)
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /admin/dialog [get]
func (c *DialogController) List(context *gin.Context) {
	var (
		waitCustomer   = []WaitCustomer{}
		kfId, _        = context.Get("KFID")
		roomCollection = model.Db.C("room")
	)
	query := []bson.M{
		{
			"$match": bson.M{"room_kf.kf_id": kfId, "room_messages.ack": false},
		},
		{
			"$project": bson.M{
				"room_customer": 1,
				"room_messages": bson.M{
					"$filter": bson.M{
						"input": "$room_messages",
						"as":    "room_message",
						"cond": bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{"$$room_message.oper_code", common.MessageFromCustomer}},
								{"$eq": []interface{}{"$$room_message.ack", false}},
							},
						},
					},
				},
			},
		},
	}

	err := roomCollection.Pipe(query).All(&waitCustomer)
	returnwaitCustomer := waitCustomer[:0]
	for _, waitC := range waitCustomer {
		if len(waitC.RoomMessages) > 0 {
			for _, v := range waitC.RoomMessages {
				v.CreateTime = v.CreateTime.In(common.LocalLocation)
			}
			returnwaitCustomer = append(returnwaitCustomer, waitC)
		}
	}
	count, _ := roomCollection.Find(bson.M{"room_kf.kf_id": "", "room_messages.oper_code": common.MessageFromCustomer}).Count()
	if err != nil {
		ReturnErrInfo(context, err)
	}

	context.JSON(http.StatusOK, struct {
		WaitReplyLists []WaitCustomer `json:"wait_reply_lists"`
		WaitQueueCount int            `json:"wait_queue_count"`
	}{
		WaitReplyLists: returnwaitCustomer,
		WaitQueueCount: count,
	})
}

// @Summary 待接入列表
// @Description 待接入列表
// @Tags WaitQueue
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /admin/wait_queue [get]
func (c *DialogController) Queue(context *gin.Context) {
	var (
		waitCustomer   = []WaitCustomer{}
		roomCollection = model.Db.C("room")
	)

	query := []bson.M{
		{
			"$match": bson.M{"room_kf.kf_id": "", "room_messages.oper_code": common.MessageFromCustomer},
		},
		{
			"$project": bson.M{
				"room_customer": 1,
				"room_messages": bson.M{"$slice": []interface{}{"$room_messages", 0, 2}},
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
// @Router /admin/wait_queue/access [post]
func (c *DialogController) Access(context *gin.Context) {
	var (
		err               error
		aRequest          CustomerIdsRequest
		kfModel           model.Kf
		kfId, _           = context.Get("KFID")
		roomCollection    = model.Db.C("room")
		kfCollection      = model.Db.C("kefu")
		messageCollection = model.Db.C("message")
	)

	if err = context.BindJSON(&aRequest); err != nil {
		ReturnErrInfo(context, err)
	}

	kfCollection.Find(bson.M{"id": kfId}).One(&kfModel)
	for _, v := range aRequest.CustomerIds {
		// 客服加入聊天房间
		roomKf := model.RoomKf{
			KfId:         kfModel.Id,
			KfName:       kfModel.NickName,
			KfHeadImgUrl: kfModel.HeadImgUrl,
		}
		// 更新会话信息
		if err = roomCollection.Update(bson.M{"room_customer.customer_id": v}, bson.M{"$set": bson.M{"room_kf": roomKf}}); err != nil {
			ReturnErrInfo(context, err)
		}
		// 归档历史会话
		if err = messageCollection.Update(bson.M{"customer_id": v}, bson.M{"$set": bson.M{"kf_id": roomKf.KfId}}); err != nil {
			// 暂停历史回话报错
			log.Warn(err)
		}
	}

	ReturnSuccessInfo(context)
}

// @Summary 确认已读
// @Description 确认已读
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /admin/dialog/ack [put]
func (c *DialogController) Ack(context *gin.Context) {
	var (
		aRequest       CustomerIdsRequest
		kfId, _        = context.Get("KFID")
		roomCollection = model.Db.C("room")
	)
	if bindErr := context.BindJSON(&aRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	for _, v := range aRequest.CustomerIds {
		if _, updateErr := roomCollection.UpdateWithArrayFilters(
			bson.M{"room_kf.kf_id": kfId, "room_customer.customer_id": v, "room_messages.oper_code": common.MessageFromCustomer},
			bson.M{"$set": bson.M{"room_messages.$[e].ack": true}},
			[]bson.M{{"e.oper_code": common.MessageFromCustomer}},
			true); updateErr != nil {
			log.Warn(updateErr)
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
// @Router /admin/dialog/{customerId}/{page}/{limit} [get]
func (c *DialogController) History(context *gin.Context) {
	var (
		roomHistory    RoomHistory
		customerId     = context.Param("customerId")
		strPage        = context.Param("page")
		strLimit       = context.Param("limit")
		roomCollection = model.Db.C("room")
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
			"$unwind": "$room_messages",
		},
		{
			"$sort": bson.M{"room_messages.create_time": 1},
		},
		{
			"$skip": (page - 1) * limit,
		},
		{
			"$limit": limit,
		},
		{
			"$group": bson.M{
				"_id":           "$_id",
				"room_messages": bson.M{"$push": "$room_messages"},
			},
		},
	}
	if err := roomCollection.Pipe(query).One(&roomHistory); err != nil {
		log.Warn(err)
	}

	for _, v := range roomHistory.RoomMessages {
		v.FormatterTimeLocation()
	}

	context.JSON(http.StatusOK, roomHistory)
}

// @Summary 发送消息
// @Description 发送消息
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /admin/dialog [post]
func (c *DialogController) SendMessage(context *gin.Context) {
	var (
		sendRequest        SendMessageRequest
		kfId, _            = context.Get("KFID")
		roomCollection     = model.Db.C("room")
		customerCollection = model.Db.C("customer")
	)
	if bindErr := context.Bind(&sendRequest); bindErr != nil {
		ReturnErrInfo(context, bindErr)
	}

	// 实时存储
	query := bson.M{
		"room_kf.kf_id":             kfId,
		"room_customer.customer_id": sendRequest.CustomerId,
	}
	changes := bson.M{
		"$push": bson.M{"room_messages": bson.M{"$each": []model.RoomMessage{
			{
				Id:         common.GetNewUUID(),
				Type:       sendRequest.MsgType,
				Msg:        sendRequest.Msg,
				OperCode:   common.MessageFromKf,
				CreateTime: time.Now(),
			},
		},
			"$slice": -100}},
	}
	if err := roomCollection.Update(query, changes); err != nil {
		ReturnErrInfo(context, errors.New("发送消息异常，存储异常，未发送成功"))
	}

	// 历史存储
	// 存储历史消息
	model.InsertMessage(model.Message{
		Id:         common.GetNewUUID(),
		Type:       sendRequest.MsgType,
		CustomerId: sendRequest.CustomerId,
		Msg:        sendRequest.Msg,
		OperCode:   common.MessageFromKf,
		CreateTime: time.Now(),
	})

	customer := model.Customer{}
	customerCollection.Find(bson.M{"customer_id": sendRequest.CustomerId}).One(&customer)
	if customer.CustomerSourceType == common.FromWeixin {
		msgResponse, err := c.wxContext.GetKf().Send(kf.KfSendMsgRequest{
			ToUser:  sendRequest.CustomerId,
			MsgType: sendRequest.MsgType,
			Text: message.Text{
				Content: strings.Replace(sendRequest.Msg, "<br>", "\n", -1),
			},
		})
		ReturnErrInfo(context, err)
		log.Printf("客服[%s]发送信息：%s \n", kfId, sendRequest.Msg)
		if msgResponse.ErrCode == 0 {
			ReturnSuccessInfo(context)
		} else {
			ReturnErrInfo(context, errors.New("发送消息失败"))
		}
	} else {
		ReturnSuccessInfo(context)
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
	RoomMessages []*model.RoomMessage `bson:"room_messages" json:"room_messages"`
}
