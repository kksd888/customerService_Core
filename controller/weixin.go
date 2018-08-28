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
		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)

		log.Printf("用户[%s]发来信息：%s \n", msg.FromUserName, text.Content)

		room, isNew := logic.InitRoom(msg.FromUserName)
		room.AddMessage(text.Content)

		if isNew {
			userInfo, err := c.wxContext.GetUser().GetUserInfo(msg.FromUserName)
			if err != nil {
				log.Fatalf("WeiXinController.wxContext.GetUser().GetUserInfo() is err：%v", err.Error())
			}

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
