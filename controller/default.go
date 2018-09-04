// 健康检查

package controller

import (
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
		kf             = model.Kf{}
		waitCustomer   = []WaitCustomer{}
		onlineCustomer = []OnlineCustomer{}
		kfId, _        = context.Get("KFID")
		kfCollection   = c.db.C("kf")
		roomCollection = c.db.C("room")
	)

	kfCollection.Find(bson.M{"id": kfId}).One(&kf)
	roomCollection.Find(bson.M{"room_kf.kf_id": kfId}).All(&onlineCustomer)
	roomCollection.Find(bson.M{"room_kf.kf_id": ""}).All(&waitCustomer)

	context.JSON(http.StatusOK, InitResponse{
		Mine: InitMine{
			Id:         kf.Id,
			UserName:   kf.NickName,
			HeadImgUrl: kf.HeadImgUrl,
			Status:     kf.Status,
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
	Mine           InitMine         `bson:"mine" json:"mine"`
	OnlineCustomer []OnlineCustomer `bson:"online_customer" json:"online_customer"`
	WaitCustomer   []WaitCustomer   `bson:"wait_customer" json:"wait_customer"`
}
type InitMine struct {
	Id         string `bson:"id" json:"id"`
	UserName   string `bson:"user_name" json:"user_name"`
	HeadImgUrl string `bson:"head_img_url" json:"head_img_url"`
	Status     bool   `json:"is_online" json:"status"`
}
type CustomerInfo struct {
	CustomerId         string `bson:"customer_id" json:"customer_id"`
	CustomerNickName   string `bson:"customer_nick_name" json:"customer_nick_name"`
	CustomerHeadImgUrl string `bson:"customer_head_img_url" json:"customer_head_img_url"`
}
type RoomMessage struct {
	Id         string    `bson:"id" json:"id"`
	Type       string    `bson:"type" json:"message_type"`
	Msg        string    `bson:"msg" json:"message_content"`
	OperCode   int       `bson:"oper_code" json:"message_oper_code"`
	Ack        bool      `bson:"ack" json:"message_ack"`
	CreateTime time.Time `bson:"create_time" json:"create_time"`
}
type OnlineCustomer struct {
	RoomCustomer CustomerInfo  `bson:"room_customer" json:"room_customer"`
	RoomMessages []RoomMessage `bson:"room_messages" json:"room_messages"`
}
type WaitCustomer struct {
	RoomCustomer CustomerInfo  `bson:"room_customer" json:"room_customer"`
	RoomMessages []RoomMessage `bson:"room_messages" json:"room_messages"`
	//PreviousKf         WaitQueuePreviousKf
}
