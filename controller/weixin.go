package controller

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/cache"
	"git.jsjit.cn/customerService/customerService_Core/wechat/message"
	"github.com/gin-gonic/gin"
	"log"
)

type WeiXinController struct {
}

func InitWeiXin() *WeiXinController {
	return &WeiXinController{}
}

var (
	Wc         *wechat.Wechat
	WxMsgQueue = make(chan *message.MixMessage, 10)
)

func init() {
	redis := cache.NewRedis(&cache.RedisOpts{
		Host: "localhost:32768",
	})

	//配置微信参数
	config := &wechat.Config{
		AppID:          "wx6cfceff5167a6007",
		AppSecret:      "1c1a365155e23b491f4878afbb87b918",
		Token:          "1603411701",
		EncodingAESKey: "fTrvMnac80fBHFP63KTLFZAhfxdSq7c126yftPw3HO1",
		Cache:          redis,
	}
	Wc = wechat.NewWechat(config)
}

// 微信通信接口
func (c *WeiXinController) Listen(context *gin.Context) {

	wcServer := Wc.GetServer(context.Request, context.Writer)

	//设置接收消息的处理方法
	wcServer.SetMessageHandler(func(msg message.MixMessage) *message.Reply {
		WxMsgQueue <- &msg

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
