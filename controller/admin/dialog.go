// 对话管理

package admin

import (
	"customerService_Core/common"
	"customerService_Core/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/li-keli/go-tool/util/mongo_util"
	"github.com/li-keli/go-tool/wechat"
	"github.com/li-keli/go-tool/wechat/kf"
	"github.com/li-keli/go-tool/wechat/message"
	"github.com/li-keli/mgo/bson"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

type DialogController struct {
	wxContext *wechat.Wechat
}

func NewDialog(wxContext *wechat.Wechat) *DialogController {
	return &DialogController{wxContext: wxContext}
}

// @Summary 会话确认应答
// @Description 会话确认应答
// @Tags WaitQueue
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /admin/wait_queue/access [post]
func (c *DialogController) Access(context *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		err               error
		aRequest          CustomerIdsRequest
		kfModel           model.Kf
		kfId              = context.GetString("KFID")
		roomCollection    = session.DB(common.AppConfig.DbName).C("room")
		kfCollection      = session.DB(common.AppConfig.DbName).C("kefu")
		messageCollection = session.DB(common.AppConfig.DbName).C("message")
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

	// websocket 通知给客服，同时广播此用户已被接入
	for _, customerId := range aRequest.CustomerIds {
		SendMsgToOnlineKf(kfId, WebSocketConnModel{Type: 1, Body: customerId})
		SendMsgRadio(WebSocketConnModel{Type: 2, Body: customerId})
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
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		aRequest       CustomerIdsRequest
		kfId           = context.GetString("KFID")
		roomCollection = session.DB(common.AppConfig.DbName).C("room")
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

// @Summary 发送消息
// @Description 发送消息
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /admin/dialog [post]
func (c *DialogController) SendMessage(context *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		sendRequest        SendMessageRequest
		kfId, _            = context.Get("KFID")
		roomCollection     = session.DB(common.AppConfig.DbName).C("room")
		customerCollection = session.DB(common.AppConfig.DbName).C("customer")
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
				KfId:       kfId.(string),
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
		KfId:       kfId.(string),
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
