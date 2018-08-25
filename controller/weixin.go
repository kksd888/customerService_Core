package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"github.com/gin-gonic/gin"
)

type WeiXinController struct {
}

func InitWeiXin() *WeiXinController {
	return &WeiXinController{}
}

// 微信通信接口
func (c *WeiXinController) Listen(context *gin.Context) {
	var (
		w = context.Writer
		r = context.Request
	)

	token := "1603411701"                        // 微信公众平台的Token
	appid := "wx6cfceff5167a6007"                // 微信公众平台的AppID
	secret := "1c1a365155e23b491f4878afbb87b918" // 微信公众平台的AppSecret
	mp := wechat.New(token, appid, secret)

	// 检查请求是否有效
	// 仅主动发送消息时不用检查
	if !mp.Request.IsValid(w, r) {
		return
	}
	// 判断消息类型
	if mp.Request.MsgType == wechat.MsgTypeText {
		// 回复消息
		mp.ReplyTextMsg(w, "Hello, 世界")
	}
}
