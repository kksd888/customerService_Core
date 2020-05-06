package admin

import (
	"customerService_Core/common"
	"github.com/li-keli/go-tool/util/mongo_util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/li-keli/mgo/bson"
)

type StatisticsController struct {
}

func NewStatistics() *StatisticsController {
	return &StatisticsController{}
}

// 统计查询
func (c *StatisticsController) Statistics(ctx *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		input = struct {
			StartTime time.Time `json:"StartTime" binding:"required"`
			EndTime   time.Time `json:"EndTime" binding:"required"`
		}{}
		messageCollection = session.DB(common.AppConfig.DbName).C("message")
		kfCollection      = session.DB(common.AppConfig.DbName).C("kefu")

		allKf []StatisticsKf
	)

	if err := ctx.BindJSON(&input); err != nil {
		ReturnErrInfo(ctx, err)
	}

	_ = kfCollection.Find(nil).All(&allKf)

	for k, v := range allKf {
		var result []string
		_ = messageCollection.Find(bson.M{"kf_id": v.Id, "create_time": bson.M{"$gte": input.StartTime, "$lt": input.EndTime}}).Distinct("customer_id", &result)
		cc, _ := messageCollection.Find(bson.M{"kf_id": v.Id, "create_time": bson.M{"$gte": input.StartTime, "$lt": input.EndTime}}).Count()

		allKf[k].SendCount = cc
		allKf[k].ReceptionCount = len(result)
	}

	ctx.JSON(http.StatusOK, allKf)
}

type StatisticsKf struct {
	Id             string `json:"-" bson:"id"`
	JobNum         string `json:"JobNum" bson:"job_num"`
	NickName       string `json:"NickName" bson:"nick_name"`
	HeadImgUrl     string `json:"HeadImgUrl" bson:"head_img_url"`
	GroupName      string `json:"GroupName" bson:"group_name"`
	SendCount      int    `json:"SendCount"`
	ReceptionCount int    `json:"ReceptionCount"`
}
type StatisticsDistinctMsg struct {
	Id string `bson:"id"`
}
