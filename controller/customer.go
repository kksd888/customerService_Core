// 访客操作

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
)

type CustomerController struct {
	wxContext *wechat.Wechat
	rooms     map[string]*logic.Room
}

func InitCustomer(wxContext *wechat.Wechat, rooms map[string]*logic.Room) *CustomerController {
	return &CustomerController{wxContext, rooms}
}
