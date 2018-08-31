// 访客操作

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
)

type CustomerController struct {
	db        *model.MongoDb
	wxContext *wechat.Wechat
}

func InitCustomer(wxContext *wechat.Wechat, _db *model.MongoDb) *CustomerController {
	return &CustomerController{wxContext: wxContext, db: _db}
}
