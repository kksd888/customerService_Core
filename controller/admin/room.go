package admin

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/li-keli/go-tool/util/db_util"
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
	session := db_util.MongoDbSession.Copy()
	defer session.Close()

	var (
		changeKfSruct = struct {
			GroupName string `json:"group_name"`
			RoomId    string `json:"room_id"`
		}{}
		kfCollection = session.DB(common.AppConfig.DbName).C("kefu")
		kfId, _      = context.Get("KFID")
		//kfId = "90a43c5cdbd34e90ae0f23af90698d86"
	)

	if err := context.Bind(&changeKfSruct); err != nil {
		ReturnErrInfo(context, errors.New(fmt.Sprintf("切换客服参数错误：%s", err.Error())))
	}

	if changeKfSruct.RoomId == "" || changeKfSruct.GroupName == "" {
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
			kfId := kfOnline[seed.Intn(len(kfOnline))].Id
			mesCollection := session.DB(common.AppConfig.DbName).C("message")
			if e := mesCollection.Update(bson.M{"id": changeKfSruct.RoomId}, bson.M{"$set": bson.M{"kf_id": kfId}}); e != nil {
				ReturnErrInfo(context, err)
			}
		} else {
			ReturnErrInfo(context, "未查到在线客服!")
		}

	}
	context.JSON(http.StatusOK, LoginInResponse{
		Authentication: kfId,
	})

}
