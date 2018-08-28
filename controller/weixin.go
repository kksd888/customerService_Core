package controller

import (
	"fmt"
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
		room.AddMessage(text.Content)

		if isNew {
			userInfo, err := c.wxContext.GetUser().GetUserInfo(msg.FromUserName)
			if err != nil {
				log.Fatalf("WeiXinController.wxContext.GetUser().GetUserInfo() is err：%v", err.Error())
			}

			// 填充附属数据
			logic.UpdateRoom(&logic.Room{
				RoomCustomer: logic.RoomCustomer{
					CustomerId:         room.CustomerId,
					CustomerNickName:   userInfo.Nickname,
					CustomerHeadImgUrl: userInfo.Headimgurl,
					CustomerMsgs:       room.CustomerMsgs,
				},
				RoomKf: logic.RoomKf{
					KfId:         room.KfId,
					KfName:       room.KfName,
					KfHeadImgUrl: room.KfHeadImgUrl,
					KfStatus:     room.KfStatus,
				},
				CreateTime: room.CreateTime,
			})

			// 客户数据持久化
			model.Customer{
				OpenId:       msg.FromUserName,
				NickName:     userInfo.Nickname,
				CustomerType: 1,
				Sex:          userInfo.Sex,
				HeadImgUrl:   userInfo.Headimgurl,
				Address:      fmt.Sprintf("%s_%s", userInfo.Province, userInfo.City),
			}.InsertOrUpdate()
		}

		//logic.PrintRoomMap()

		return
	})

	//处理消息接收以及回复
	err := wcServer.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	//发送接收成功
	wcServer.SendSuccess()
}
