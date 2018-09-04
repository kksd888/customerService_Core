package controller

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

type WeiXinController struct {
	db        *model.MongoDb
	wxContext *wechat.Wechat
}

func InitWeiXin(wxContext *wechat.Wechat, _db *model.MongoDb) *WeiXinController {
	return &WeiXinController{wxContext: wxContext, db: _db}
}

// 微信通信接口
func (c *WeiXinController) Listen(context *gin.Context) {
	wcServer := c.wxContext.GetServer(context.Request, context.Writer)

	//设置接收消息的处理方法
	wcServer.SetMessageHandler(func(msg message.MixMessage) (reply *message.Reply) {
		/*
			A 24小时新接入客户：
				1. 注册分配聊天房间
				2. 存储新客户、离线留言信息数据
			B 已在线的客户：
				1. 检索已分配的房间
				2. 存储聊天数据
		*/
		var (
			msgType  = string(msg.MsgType) // 消息类型
			MediaUrl = ""                  // 多媒体地址
			msgText  = ""                  // 文本内容
		)
		switch msg.MsgType {
		case message.MsgTypeText:
			msgText = message.NewText(msg.Content).Content
		case message.MsgTypeImage:
			MediaUrl = msg.PicURL
		case message.MsgTypeVoice:
			msgText = msg.Recognition
			material := c.wxContext.GetMaterial()
			if mediaURL, err := material.GetMediaURL(msg.MediaID); err != nil {
				log.Printf("material.MsgTypeVoice is err: %#v", err)
			} else {
				MediaUrl = mediaURL
			}
		case message.MsgTypeVideo:
			material := c.wxContext.GetMaterial()
			if mediaURL, err := material.GetMediaURL(msg.MediaID); err != nil {
				log.Printf("material.MsgTypeVideo is err: %#v", err)
			} else {
				MediaUrl = mediaURL
			}
		case message.MsgTypeShortVideo:
			material := c.wxContext.GetMaterial()
			if mediaURL, err := material.GetMediaURL(msg.MediaID); err != nil {
				log.Printf("material.MsgTypeShortVideo is err: %#v", err)
			} else {
				MediaUrl = mediaURL
			}
		}

		log.Printf("用户[%s]发来信息：[%s] %s \n", msg.FromUserName, msgType, msgText)

		roomCollection := c.db.C("room")
		customerCollection := c.db.C("customer")
		count, _ := roomCollection.Find(bson.M{"room_customer.customer_id": msg.FromUserName}).Count()

		if count == 0 {
			// 新接入
			userInfo, err := c.wxContext.GetUser().GetUserInfo(msg.FromUserName)
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
			})

			//if _, isOk := logic.GetOnlineKf(); !isOk {
			//	return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(common.KF_REPLY)}
			//}

			roomCollection.Insert(&model.Room{
				RoomCustomer: model.RoomCustomer{
					CustomerId:         msg.FromUserName,
					CustomerNickName:   userInfo.Nickname,
					CustomerHeadImgUrl: userInfo.Headimgurl,
				},
				RoomMessages: []model.RoomMessage{
					{
						Id:         common.GetNewUUID(),
						Type:       msgType,
						Msg:        msgText,
						MediaUrl:   MediaUrl,
						OperCode:   common.MessageFromCustomer,
						CreateTime: time.Now(),
					},
				},
				CreateTime: time.Now(),
			})
		} else {
			// 会话中的客户
			roomCollection.Update(bson.M{"room_customer.customer_id": msg.FromUserName},
				bson.M{"$push": bson.M{"room_messages": model.RoomMessage{
					Id:         common.GetNewUUID(),
					Type:       msgType,
					Msg:        msgText,
					MediaUrl:   MediaUrl,
					OperCode:   common.MessageFromCustomer,
					CreateTime: time.Now(),
				}}})
		}

		return
	})

	//处理消息接收以及回复
	err := wcServer.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	//发送接收成功
	wcServer.Send()
}
