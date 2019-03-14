package admin

import (
	"customerService_Core/common"
	"customerService_Core/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/li-keli/go-tool/util/mongo_util"
	"github.com/li-keli/mgo/bson"
	"github.com/pkg/errors"
	"math/rand"
	"net/http"
	"time"
)

type RoomController struct {
}

func NewRoom() *RoomController {
	return &RoomController{}
}

// 房间客服变更
func (c *RoomController) ChangeKf(context *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		changeKfSruct = struct {
			GroupName  string `json:"group_name"`
			CustomerId string `json:"customer_id"`
		}{}
		kfCollection = session.DB(common.AppConfig.DbName).C("kefu")
		kfId, _      = context.Get("KFID")
		changeKfId   = ""
	)

	if err := context.Bind(&changeKfSruct); err != nil {
		ReturnErrInfo(context, errors.New(fmt.Sprintf("切换客服参数错误：%s", err.Error())))
	}

	if changeKfSruct.CustomerId == "" || changeKfSruct.GroupName == "" {
		ReturnErrInfo(context, "切换客服参数错误")
	}

	kfOnline := []model.Kf{}

	if err := kfCollection.Find(bson.M{
		"group_name": changeKfSruct.GroupName,
		"is_online":  true,
		"id":         bson.M{"$ne": kfId},
	}).All(&kfOnline); err != nil {
		ReturnErrInfo(context, err)
	} else {
		if len(kfOnline) > 0 {
			seed := rand.New(rand.NewSource(time.Now().UnixNano()))
			rangNo := seed.Intn(len(kfOnline))
			changeKfId = kfOnline[rangNo].Id
			mesCollection := session.DB(common.AppConfig.DbName).C("room")
			if e := mesCollection.Update(bson.M{"room_customer.customer_id": changeKfSruct.CustomerId}, bson.M{"$set": bson.M{"room_kf.kf_id": changeKfId, "room_kf.kf_name": kfOnline[rangNo].NickName, "room_kf.kf_head_img_url": kfOnline[rangNo].HeadImgUrl, "room_kf.kf_status": kfOnline[rangNo].IsOnline}}); e != nil {
				ReturnErrInfo(context, err)
			}
		} else {
			ReturnErrInfo(context, "未查到在线客服!")
		}

	}
	context.JSON(http.StatusOK, LoginInResponse{
		Authentication: changeKfId,
	})

}
