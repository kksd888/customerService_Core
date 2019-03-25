package admin

import (
	"customerService_Core/common"
	"customerService_Core/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/li-keli/go-tool/util/mongo_util"
	"github.com/li-keli/mgo/bson"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
)

type RoomController struct {
}

func NewRoom() *RoomController {
	return &RoomController{}
}

// 房间客服变更
func (c *RoomController) Transfer(context *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		changeKfSruct = struct {
			CustomerId   string `json:"customer_id" binding:"required"`
			TransferKfId string `json:"transfer_kf_id" binding:"required"`
		}{}
		kfCollection  = session.DB(common.AppConfig.DbName).C("kefu")
		mesCollection = session.DB(common.AppConfig.DbName).C("room")
		//kfId, _       = context.Get("KFID")
		kfModel model.Kf
	)

	if err := context.Bind(&changeKfSruct); err != nil {
		ReturnErrInfo(context, errors.New(fmt.Sprintf("切换客服参数错误：%s", err.Error())))
	}

	if err := kfCollection.Find(bson.M{"id": changeKfSruct.TransferKfId, "is_online": true}).One(&kfModel); err != nil {
		ReturnErrInfo(context, err)
	}

	logrus.Info(changeKfSruct.CustomerId)
	logrus.Info(kfModel)

	if err := mesCollection.Update(
		bson.M{"room_customer.customer_id": changeKfSruct.CustomerId},
		bson.M{"$set": bson.M{
			"room_kf.kf_id":           kfModel.Id,
			"room_kf.kf_name":         kfModel.NickName,
			"room_kf.kf_head_img_url": kfModel.HeadImgUrl,
			"room_kf.kf_status":       true,
		}}); err != nil {
		ReturnErrInfo(context, err)
	}

	context.JSON(http.StatusOK, nil)
}
