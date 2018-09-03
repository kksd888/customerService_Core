// 健康检查

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
)

type DefaultController struct {
	db *model.MongoDb
}

func InitHealth(_db *model.MongoDb) *DefaultController {
	return &DefaultController{db: _db}
}

// @Summary 健康检查˚
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
	// 获取访问客服信息
	var (
		onlineCustomer []OnlineCustomer
		waitCustomer   []WaitCustomer
		kfId, _        = context.Get("KFID")
		kf             = model.Kf{}
	)

	kfCollection := c.db.C("kf")
	roomCollection := c.db.C("room")

	kfCollection.Find(bson.M{"id": kfId}).One(&kf)
	roomCollection.Find(bson.M{"roomkf.kfid": kfId}).All(&onlineCustomer)
	roomCollection.Find(bson.M{"roomkf.kfid": ""}).All(&waitCustomer)

	context.JSON(http.StatusOK, InitResponse{
		Mine: InitMine{
			Id:         kf.Id,
			UserName:   kf.NickName,
			HeadImgUrl: kf.HeadImgUrl,
			Status:     common.KF_ONLINE,
		},
		OnlineCustomer: onlineCustomer,
		WaitCustomer:   waitCustomer,
	})
}

// 异常返回
func ReturnErrInfo(context *gin.Context, err interface{}) {
	if err != nil {
		log.Printf("发生异常：%#v", err)
		context.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  err.(error).Error(),
		})
		panic(err)
	}
}

// 成功返回
func ReturnSuccessInfo(context *gin.Context) {
	context.JSON(http.StatusInternalServerError, gin.H{
		"code": http.StatusOK,
		"msg":  "ok",
	})
}

type InitResponse struct {
	Mine           InitMine         `json:"mine"`
	OnlineCustomer []OnlineCustomer `json:"online_customer"`
	WaitCustomer   []WaitCustomer   `json:"wait_customer"`
}
type InitMine struct {
	Id         string `json:"id"`
	UserName   string `json:"user_name"`
	HeadImgUrl string `json:"head_img_url"`
	Status     int    `json:"status"`
}
type CustomerInfo struct {
	CustomerId         string `json:"customer_id"`
	CustomerNickName   string `json:"customer_nick_name"`
	CustomerHeadImgUrl string `json:"customer_head_img_url"`
}
type RoomMessage struct {
	Id         string    `json:"id"`
	Type       string    `json:"message_type"`
	Msg        string    `json:"message_content"`
	OperCode   int       `json:"message_oper_code"`
	Ack        bool      `json:"message_ack"`
	CteateTime time.Time `json:"message_cteate_time"`
}
type OnlineCustomer struct {
	RoomCustomer CustomerInfo  `json:"room_customer"`
	RoomMessages []RoomMessage `json:"room_messages"`
}
type WaitCustomer struct {
	RoomCustomer CustomerInfo  `json:"room_customer"`
	RoomMessages []RoomMessage `json:"room_messages"`
	//PreviousKf         WaitQueuePreviousKf
}
