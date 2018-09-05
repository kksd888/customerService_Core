// 访客操作

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/wechat"
)

type CustomerController struct {
	wxContext *wechat.Wechat
}

func InitCustomer(wxContext *wechat.Wechat) *CustomerController {
	return &CustomerController{wxContext: wxContext}
}
