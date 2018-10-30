package open

import (
	"encoding/base64"
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type DialogController struct {
	aiModule *handle.AiSemantic
}

func NewDialog(aiModule *handle.AiSemantic) *DialogController {
	return &DialogController{aiModule: aiModule}
}

// 获取历史记录
// /v1/app/dialog/history
func (dialog *DialogController) History(ctx *gin.Context) {
	var (
		customerId, _ = ctx.Get("CID")
		output        []MessageModel
		dbResult      []struct {
			RoomMessages MessageModel `bson:"room_messages"`
		}
		roomCollection = model.Db.C("room")
	)

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
			"$project": bson.M{"_id": 0, "room_messages": 1},
		},
	}
	if err := roomCollection.Pipe(query).All(&dbResult); err != nil {
		common.ReturnErr(ctx, err)
	}
	for _, v := range dbResult {
		v.RoomMessages.CreateTime2Timestamp()
		output = append(output, v.RoomMessages)
	}

	// 将所有未读的客服消息进行已读操作
	if _, updateErr := roomCollection.UpdateWithArrayFilters(
		bson.M{"room_customer.customer_id": customerId},
		bson.M{"$set": bson.M{"room_messages.$[e].ack": true}},
		[]bson.M{{"e.oper_code": common.MessageFromKf}},
		true); updateErr != nil {
		log.Warn(updateErr)
	}

	// 清除系统提示消息
	roomCollection.Update(bson.M{"room_customer.customer_id": customerId},
		bson.M{"$pull": bson.M{"room_messages": bson.M{"oper_code": common.MessageFromSys}}})

	common.ReturnSuccess(ctx, output)
}

// 获取新消息
// /v1/app/dialog
func (dialog *DialogController) Get(ctx *gin.Context) {
	var (
		customerId, _ = ctx.Get("CID")
		output        []MessageModel
		dbResult      []struct {
			RoomMessages MessageModel `bson:"room_messages"`
		}

		roomCollection = model.Db.C("room")
	)
	query := []bson.M{
		{
			"$match": bson.M{"room_customer.customer_id": customerId, "room_messages.oper_code": common.MessageFromKf, "room_messages.ack": false},
		},
		{
			"$project": bson.M{
				"_id": 0,
				"room_messages": bson.M{
					"$filter": bson.M{
						"input": "$room_messages",
						"as":    "room_message",
						"cond": bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{"$$room_message.oper_code", common.MessageFromKf}},
								{"$eq": []interface{}{"$$room_message.ack", false}},
							},
						},
					},
				},
			},
		},
		{
			"$unwind": "$room_messages",
		},
	}

	if err := roomCollection.Pipe(query).All(&dbResult); err != nil {
		common.ReturnErr(ctx, err)
	}
	for _, v := range dbResult {
		v.RoomMessages.CreateTime2Timestamp()
		output = append(output, v.RoomMessages)
	}

	// 确认已读的消息
	if _, updateErr := roomCollection.UpdateWithArrayFilters(
		bson.M{"room_customer.customer_id": customerId},
		bson.M{"$set": bson.M{"room_messages.$[e].ack": true}},
		[]bson.M{{"e.oper_code": common.MessageFromKf}},
		true); updateErr != nil {
		log.Warn(updateErr)
	}

	common.ReturnSuccess(ctx, output)
}

// 发送消息
// /v1/app/dialog
func (dialog *DialogController) Create(ctx *gin.Context) {
	var (
		err           error
		uploadAddress string
		fileAddress   string
		customerId    = ctx.GetString("CID")
		input         struct {
			Msg           string         `json:"msg"`                     // 文本消息
			Type          common.MsgType `json:"type" binding:"required"` // 多媒体类型
			ExtensionName string         `json:"extension_name"`          // 媒体扩展名
			MediaBase64   string         `json:"media_base64"`            // 多媒体base64
		}
	)

	if err = ctx.BindJSON(&input); err != nil {
		common.ReturnErrCode(ctx, common.ParameterBad, err)
	}

	if input.Type == common.MsgTypeText && input.Msg == "" {
		common.ReturnErrCode(ctx, common.ParameterBad, errors.New("不能发送空的文本"))
	}

	// 解析上传的文件
	if input.MediaBase64 != "" {
		if input.ExtensionName == "" {
			common.ReturnErrCode(ctx, common.ParameterBad, errors.New("上传文件缺少扩展名"))
		}
		switch input.Type {
		case common.MsgTypeText, common.MsgTypeImage, common.MsgTypeVoice:
			// 存储位置格式：upload/多媒体类型/日期
			uploadAddress = fmt.Sprintf("./upload/%s/%s", string(input.Type), time.Now().Format("2006-01"))
			if _, err = os.Stat(uploadAddress); err != nil {
			}
			if os.IsNotExist(err) {
				err = os.MkdirAll(uploadAddress, os.ModePerm)
				common.ReturnErr(ctx, err)
			}
			bytes, _ := base64.StdEncoding.DecodeString(input.MediaBase64)
			fileAddress = fmt.Sprintf("%s/%s.%s", uploadAddress, common.GetNewUUID(), input.ExtensionName)
			ioutil.WriteFile(fileAddress, bytes, 0666)
		default:
			common.ReturnErrCode(ctx, common.ParameterBad, errors.New("未知的多媒体类型，请检查Type值"))
		}
	}

	// 发送给客服调度模块
	var mediaUrl = ""
	if fileAddress != "" {
		mediaUrl = "http://kf.api.7u1.cn/" + strings.TrimLeft(fileAddress, "./") // 去掉虚拟路径的前缀
	}

	sendReply := dialog.send(SendModel{
		FromUserName: customerId,
		Msg:          input.Msg,
		MsgType:      string(input.Type),
		MediaUrl:     mediaUrl,
	})

	if sendReply != "" {
		common.ReturnSuccess(ctx, struct {
			Reply string `json:"reply"`
		}{sendReply})
		return
	}

	common.ReturnSuccess(ctx, gin.H{})
}

// 消息模型
type MessageModel struct {
	Id         string         `json:"id" bson:"id"`               // 消息唯一编号
	Type       common.MsgType `json:"type" bson:"type"`           // 消息类型
	MediaUrl   string         `json:"media_url" bson:"media_url"` // 多媒体消息地址
	Msg        string         `json:"msg" bson:"msg"`             // 文本消息正文
	Ack        bool           `json:"ack" bson:"ack"`             // 是否已读
	OperCode   int            `json:"oper_code" bson:"oper_code"` // 操作码
	Timestamp  int64          `json:"timestamp"`                  // 创建时间
	CreateTime time.Time      `json:"-" bson:"create_time"`
}

// 创建时间转时间戳
func (m *MessageModel) CreateTime2Timestamp() {
	m.Timestamp = m.CreateTime.Unix()
}

type SendModel struct {
	FromUserName string `json:"from_user_name"` // 发送者
	Msg          string `json:"msg"`            // 文本消息
	MsgType      string `json:"type"`           // 多媒体类型
	MediaUrl     string `json:"media_url"`      // 多媒体地址
}

// 监听移动模块发送过来的消息
func (dialog *DialogController) send(msg SendModel) string {
	var (
		roomCollection = model.Db.C("room")
		aiDialogue     = "" // AI答复

	)
	// 小金尝试回答
	if msg.Msg != "" {
		aiDialogue = dialog.aiModule.Dialogue(msg.Msg, msg.FromUserName)
	}

	if aiDialogue != "" {
		log.Printf("用户[%s]发来信息：[%s] %s；小金推荐回复：%s \n", msg.FromUserName, msg.MsgType, msg.Msg, aiDialogue)
	} else {
		log.Printf("用户[%s]发来信息：[%s] %s \n", msg.FromUserName, msg.MsgType, msg.Msg)
	}

	var room = model.Room{}
	roomCollection.Find(bson.M{"room_customer.customer_id": msg.FromUserName}).One(&room)

	if room.RoomCustomer.CustomerId == "" {
		// 此处不应该存在新接入的用户
	} else {
		var (
			kefuColection = model.Db.C("kefu")
			kefuModel     = model.Kf{}
		)
		kefuColection.Find(bson.M{"id": room.RoomKf.KfId}).One(&kefuModel)
		if kefuModel.Id != "" && kefuModel.IsOnline == false {
			// 若接待的客服已经下线，则将用户重新放入待接入
			roomCollection.Update(
				bson.M{"room_customer.customer_id": msg.FromUserName},
				bson.M{"$set": bson.M{"room_kf": &model.RoomKf{}}})
		}

		// 实时会话数据更新
		query := bson.M{
			"room_customer.customer_id": msg.FromUserName,
		}
		changes := bson.M{
			"$push": bson.M{"room_messages": bson.M{"$each": []model.RoomMessage{
				{
					Id:         common.GetNewUUID(),
					Type:       string(msg.MsgType),
					Msg:        msg.Msg,
					AiMsg:      aiDialogue,
					MediaUrl:   msg.MediaUrl,
					OperCode:   common.MessageFromCustomer,
					CreateTime: time.Now(),
				},
			},
				"$slice": -100}},
		}
		if err := roomCollection.Update(query, changes); err != nil {
			log.Printf("实时房型数据更新异常：%s", err.Error())
		}
	}

	// 存储历史消息
	model.InsertMessage(model.Message{
		Id:         common.GetNewUUID(),
		Type:       string(msg.MsgType),
		CustomerId: msg.FromUserName,
		KfId:       room.RoomKf.KfId,
		Msg:        msg.Msg,
		AiMsg:      aiDialogue,
		MediaUrl:   msg.MediaUrl,
		OperCode:   common.MessageFromCustomer,
		CreateTime: time.Now(),
	})

	online, _ := model.Kf{}.QueryOnlines()
	if len(online) == 0 {
		return common.KF_REPLY
	} else {
		return ""
	}
}
