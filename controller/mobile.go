package controller

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type MobileController struct {
	wxContext *wechat.Wechat
	aiModule  *handle.AiSemantic
}

func NewMobile(wxContext *wechat.Wechat, aiModule *handle.AiSemantic) *MobileController {
	return &MobileController{wxContext: wxContext, aiModule: aiModule}
}

// 监听移动模块发送过来的消息
func (mobile *MobileController) Listen(ctx *gin.Context) {
	var (
		msg struct {
			FromUserName string          `json:"from_user_name"` // 发送者
			Msg          string          `json:"msg"`            // 文本消息
			MsgType      message.MsgType `json:"type"`           // 多媒体类型
			MediaUrl     string          `json:"media_url"`      // 多媒体地址
		}

		roomCollection     = model.Db.C("room")
		customerCollection = model.Db.C("customer")
		aiDialogue         = "" // AI答复

	)
	if bindErr := ctx.Bind(&msg); bindErr != nil {
		ReturnErrInfo(ctx, bindErr)
	}

	// 小金尝试回答
	if msg.Msg != "" {
		aiDialogue = mobile.aiModule.Dialogue(msg.Msg)
	}

	if aiDialogue != "" {
		log.Printf("用户[%s]发来信息：[%s] %s；小金推荐回复：%s \n", msg.FromUserName, msg.MsgType, msg.Msg, aiDialogue)
	} else {
		log.Printf("用户[%s]发来信息：[%s] %s \n", msg.FromUserName, msg.MsgType, msg.Msg)
	}

	var room = model.Room{}
	roomCollection.Find(bson.M{"room_customer.customer_id": msg.FromUserName}).One(&room)

	if room.RoomCustomer.CustomerId == "" {
		// 新接入
		userInfo, err := mobile.wxContext.GetUser().GetUserInfo(msg.FromUserName)
		if err != nil {
			log.Printf("WeiXinController.wxContext.GetUser().GetUserInfo() is err：%v", err.Error())
		}

		// 客户数据持久化
		customerCollection.Insert(&model.Customer{
			CustomerId:   msg.FromUserName,
			NickName:     userInfo.Nickname,
			CustomerType: common.NormalCustomer,
			Sex:          userInfo.Sex,
			HeadImgUrl:   userInfo.Headimgurl,
			Address:      fmt.Sprintf("%s_%s", userInfo.Province, userInfo.City),
			CreateTime:   time.Now(),
			UpdateTime:   time.Now(),
		})

		// 实时会话数据更新
		roomCollection.Insert(&model.Room{
			RoomCustomer: model.RoomCustomer{
				CustomerId:         msg.FromUserName,
				CustomerNickName:   userInfo.Nickname,
				CustomerHeadImgUrl: userInfo.Headimgurl,
			},
			RoomMessages: []model.RoomMessage{
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
			CreateTime: time.Now(),
		})
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
		ctx.JSON(http.StatusOK, gin.H{
			"reply": common.KF_REPLY,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"reply": "",
	})
}
