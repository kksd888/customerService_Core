package controller

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/logic"
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

//func listenSendSqueue() {
//	for {
//		wxMsg := <-WxSend
//		kf := WxContext.GetKf()
//		if msgResponse, err := kf.SendTextMsg(wxMsg.ToUser, kf.Context); err != nil {
//			log.Printf("%#v", msgResponse)
//		}
//	}
//}

// 微信通信接口
func (c *WeiXinController) Listen(context *gin.Context) {

	wcServer := c.wxContext.GetServer(context.Request, context.Writer)

	//设置接收消息的处理方法
	wcServer.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		room := logic.InitRoom(msg.FromUserName)
		room.Register()

		//回复消息：演示回复用户发送的消息
		text := message.NewText(msg.Content)
		log.Printf("用户[%s]发来信息%s", msg.FromUserName, text.Content)
		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})

	//处理消息接收以及回复
	err := wcServer.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}

	//发送回复的消息
	wcServer.Send()
}
