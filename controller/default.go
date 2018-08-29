// 健康检查

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type DefaultController struct {
}

func InitHealth() *DefaultController {
	return &DefaultController{}
}

// @Summary 健康检查
// @Description 应用程序健康检查接口
// @Tags Default
// @Accept json
// @Produce json
// @Success 200 {string} json ""
// @Router /v1/health [get]
func (c *DefaultController) Health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now().Format("2006-01-02 15:04:05")})
}

// @Summary 系统初始化
// @Description 在线客服系统进行初始化
// @Tags Default
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/init [get]
func (c *DefaultController) Init(context *gin.Context) {
	roomKf, err := handle.AuthToken2Model(context)
	ReturnErrInfo(context, err)

	kfDb := &model.Kf{Id: roomKf.KfId}
	if err := kfDb.Get(); err != nil {
		ReturnErrInfo(context, err.Error())
	}

	messageDb := model.MessageLinkCustomer{Message: model.Message{KfId: kfDb.Id}}
	messages, err := messageDb.GetKfHistoryMsg()
	ReturnErrInfo(context, err)

	var initOnlineCustomers []InitOnlineCustomer
	var mapCus = map[int]*InitOnlineCustomer{}

	for _, singeMsg := range messages {
		if _, ok := mapCus[singeMsg.KfId]; ok {
			mapCus[singeMsg.KfId].CustomerMessages = append(mapCus[singeMsg.KfId].CustomerMessages, InitMessage{
				Id:                singeMsg.Id,
				MessageType:       singeMsg.MsgType,
				MessageContent:    singeMsg.Msg,
				MessageOperCode:   singeMsg.OperCode,
				MessageCteateTime: singeMsg.CreateTime,
			})
		} else {
			mapCus[singeMsg.KfId] = &InitOnlineCustomer{
				CustomerId:         singeMsg.CustomerId,
				CustomerNickName:   singeMsg.CustomerNickName,
				CustomerHeadImgUrl: singeMsg.CustomerHeadImgUrl,
				CustomerMessages: []InitMessage{
					{
						Id:                singeMsg.Id,
						MessageType:       singeMsg.MsgType,
						MessageContent:    singeMsg.Msg,
						MessageOperCode:   singeMsg.OperCode,
						MessageCteateTime: singeMsg.CreateTime,
					},
				},
			}
		}
	}
	for _, v := range mapCus {
		initOnlineCustomers = append(initOnlineCustomers, *v)
	}

	context.JSON(http.StatusOK, InitResponse{
		BaseResponse: BaseResponse{},
		Mine: InitMine{
			Id:         kfDb.Id,
			UserName:   kfDb.NickName,
			HeadImgUrl: kfDb.HeadImgUrl,
			Status:     common.KF_ONLINE,
		},
		InitOnlineCustomer: initOnlineCustomers,
	})
}

func ReturnErrInfo(context *gin.Context, err interface{}) {
	if err != nil {
		log.Printf("发生异常：%#v", err)
		context.JSON(http.StatusInternalServerError, BaseResponse{
			Code: 500,
			Msg:  "接口调用异常，或联系管理员",
		})
		panic(err)
	}
}

// API全局响应基础结构
type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type InitResponse struct {
	BaseResponse
	Mine               InitMine             `json:"mine"`
	InitOnlineCustomer []InitOnlineCustomer `json:"init_online_customer"`
}
type InitMine struct {
	Id         int    `json:"id"`
	UserName   string `json:"user_name"`
	HeadImgUrl string `json:"head_img_url"`
	Status     string `json:"status"`
}
type InitOnlineCustomer struct {
	CustomerId         int           `json:"customer_id"`
	CustomerNickName   string        `json:"customer_nick_name"`
	CustomerHeadImgUrl string        `json:"customer_head_img_url"`
	CustomerMessages   []InitMessage `json:"customer_messages"`
}
type InitMessage struct {
	Id                int       `json:"id"`
	MessageType       int       `json:"message_type"`
	MessageContent    string    `json:"message_content"`
	MessageOperCode   int       `json:"message_oper_code"`
	MessageCteateTime time.Time `json:"message_cteate_time"`
}
