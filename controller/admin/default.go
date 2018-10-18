// 健康检查

package admin

import (
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type DefaultController struct {
}

func NewHealth() *DefaultController {
	return &DefaultController{}
}

// @Summary 健康检查
// @Description 应用程序健康检查接口
// @Tags Default
// @Accept json
// @Produce json
// @Success 200 {string} json ""
// @Router /admin/health [get]
func (c *DefaultController) Health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"code": 100})
}

// @Summary 系统初始化
// @Description 在线客服系统进行初始化
// @Tags Default
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /admin/init [get]
func (c *DefaultController) Init(context *gin.Context) {
	// 获取访问客服信息
	var (
		kf             = model.Kf{}
		waitCustomer   = []WaitCustomer{}
		onlineCustomer = []OnlineCustomer{}
		kfId, _        = context.Get("KFID")
		kfCollection   = model.Db.C("kefu")
		roomCollection = model.Db.C("room")
	)

	kfCollection.Find(bson.M{"id": kfId}).One(&kf)

	// 获取聊天列表 (最多输出100条)
	roomCollection.Pipe([]bson.M{
		{
			"$match": bson.M{"room_kf.kf_id": kfId},
		},
		{
			"$project": bson.M{
				"room_customer": 1,
				"room_messages": bson.M{"$slice": []interface{}{"$room_messages", -1}},
			},
		},
		{
			"$sort": bson.M{"room_messages.create_time": -1},
		},
		{
			"$limit": 100,
		},
	}).All(&onlineCustomer)

	// 获取排队列表
	roomCollection.Pipe([]bson.M{
		{
			"$match": bson.M{"room_kf.kf_id": ""},
		},
		{
			"$project": bson.M{
				"room_customer": 1,
				"room_messages": bson.M{"$slice": []interface{}{"$room_messages", -1}},
			},
		},
	}).All(&waitCustomer)

	context.JSON(http.StatusOK, InitResponse{
		Mine: InitMine{
			Id:         kf.Id,
			UserName:   kf.NickName,
			HeadImgUrl: kf.HeadImgUrl,
			IsOnline:   kf.IsOnline,
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
	context.JSON(http.StatusOK, gin.H{
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
	IsOnline   bool   `bson:"is_online" json:"status"`
}
type CustomerInfo struct {
	CustomerId         string `bson:"customer_id" json:"customer_id"`
	CustomerNickName   string `bson:"customer_nick_name" json:"customer_nick_name"`
	CustomerHeadImgUrl string `bson:"customer_head_img_url" json:"customer_head_img_url"`
}
type OnlineCustomer struct {
	RoomCustomer CustomerInfo        `bson:"room_customer" json:"room_customer"`
	RoomMessages []model.RoomMessage `bson:"room_messages" json:"room_messages"`
}
type WaitCustomer struct {
	RoomCustomer CustomerInfo        `bson:"room_customer" json:"room_customer"`
	RoomMessages []model.RoomMessage `bson:"room_messages" json:"room_messages"`
	//PreviousKf         WaitQueuePreviousKf
}
