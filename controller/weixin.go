package controller

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"log"
)

type WeiXinController struct {
	wxContext *wechat.Wechat
	rooms     map[string]*logic.Room
}

func InitWeiXin(wxContext *wechat.Wechat, rooms map[string]*logic.Room) *WeiXinController {
	return &WeiXinController{wxContext, rooms}
}

// 微信通信接口
func (c *WeiXinController) Listen(context *gin.Context) {

	wcServer := c.wxContext.GetServer(context.Request, context.Writer)

	//设置接收消息的处理方法
	wcServer.SetMessageHandler(func(msg message.MixMessage) (reply *message.Reply) {
		text := message.NewText(msg.Content)
		log.Printf("用户[%s]发来信息：%s \n", msg.FromUserName, text.Content)

		// 通信注册
		room, isNew := logic.InitRoom(msg.FromUserName)
		//room.AddMessage(text.Content)

		// 首次访问的客户
		if isNew {
			userInfo, err := c.wxContext.GetUser().GetUserInfo(msg.FromUserName)
			if err != nil {
				log.Printf("WeiXinController.wxContext.GetUser().GetUserInfo() is err：%v", err.Error())
			}

			// 自动分配客服 (随机)
			kf, isOk := logic.GetOnlineKf()
			if !isOk {
				// 存储消息
				model.Message{
					CustomerToken: room.CustomerId,
					KfId:          kf.Id,
					KfAck:         false,
					Msg:           msg.Content,
					MsgType:       1,
					OperCode:      200,
				}.InsertOrUpdate()
				return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(common.KF_REPLY)}
			}

			// 客户数据持久化
			model.Customer{
				OpenId:       msg.FromUserName,
				NickName:     userInfo.Nickname,
				CustomerType: common.NormalCustomer,
				Sex:          userInfo.Sex,
				HeadImgUrl:   userInfo.Headimgurl,
				Address:      fmt.Sprintf("%s_%s", userInfo.Province, userInfo.City),
			}.InsertOrUpdate()

			// 存储消息
			model.Message{
				CustomerToken: room.CustomerId,
				KfId:          kf.Id,
				KfAck:         false,
				Msg:           msg.Content,
				MsgType:       1,
				OperCode:      200,
			}.InsertOrUpdate()

			// 更新Room数据
			logic.UpdateRoom(&logic.Room{
				RoomCustomer: logic.RoomCustomer{
					CustomerId:         room.CustomerId,
					CustomerNickName:   userInfo.Nickname,
					CustomerHeadImgUrl: userInfo.Headimgurl,
					CustomerMsgs:       room.CustomerMsgs,
				},
				RoomKf: logic.RoomKf{
					KfId:         kf.Id,
					KfName:       kf.NickName,
					KfHeadImgUrl: kf.HeadImgUrl,
					KfStatus:     common.KF_ONLINE,
				},
				CreateTime: room.CreateTime,
			})

		} else {
			// 未分配客服的客户
			if room.KfId == 0 {
				// 自动分配客服 (随机)
				kf, isOk := logic.GetOnlineKf()
				if !isOk {
					// 存储消息
					model.Message{
						CustomerToken: room.CustomerId,
						KfId:          kf.Id,
						KfAck:         false,
						Msg:           msg.Content,
						MsgType:       1,
						OperCode:      200,
					}.InsertOrUpdate()
					return &message.Reply{MsgType: message.MsgTypeText, MsgData: message.NewText(common.KF_REPLY)}
				}

				// 更新Room数据
				room.RoomKf = logic.RoomKf{
					KfId:         kf.Id,
					KfName:       kf.NickName,
					KfHeadImgUrl: kf.HeadImgUrl,
					KfStatus:     common.KF_ONLINE,
				}
				logic.UpdateRoom(room)

				// 存储消息
				model.Message{
					CustomerToken: room.CustomerId,
					KfId:          room.KfId,
					KfAck:         false,
					Msg:           msg.Content,
					MsgType:       1,
					OperCode:      200,
				}.InsertOrUpdate()
			} else {
				// 存储消息
				model.Message{
					CustomerToken: room.CustomerId,
					KfId:          room.KfId,
					KfAck:         false,
					Msg:           msg.Content,
					MsgType:       1,
					OperCode:      200,
				}.InsertOrUpdate()
			}
		}

		//logic.PrintRoomMap()
		log.Printf("%#v", room)

		return &message.Reply{message.MsgTypeText, nil}
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
