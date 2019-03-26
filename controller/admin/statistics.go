package admin

import (
	"customerService_Core/common"
	"github.com/li-keli/go-tool/util/mongo_util"
	"github.com/sirupsen/logrus"
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

		allKf = []StatisticsKf{}
	)

	if err := ctx.BindJSON(&input); err != nil {
		ReturnErrInfo(ctx, err)
	}

	_ = kfCollection.Find(bson.M{}).All(&allKf)

	for k, v := range allKf {
		var result = []bson.M{}
		_ = messageCollection.Find(bson.M{"kf_id": v.Id, "create_time": bson.M{"$gte": input.StartTime, "$lt": input.EndTime}}).Distinct("customer_id", &result)
		cc, _ := messageCollection.Find(bson.M{"kf_id": v.Id, "create_time": bson.M{"$gte": input.StartTime, "$lt": input.EndTime}}).Count()

		logrus.Info("distinctMsg: ", result)

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

// 统计数据
// 时间区间内，客服回复的信息总量
// db.getCollection('message').find({'oper_code':2003, 'create_time':{ "$gte" : ISODate("2018-10-01T00:00:00Z"), "$lt" : ISODate("2018-10-31T00:00:00Z")}}).count()
// 时间区间内，接待的客户数量
// db.getCollection('message').distinct('customer_id',{'oper_code':2003, 'create_time':{ "$gte" : ISODate("2018-10-01T00:00:00Z"), "$lt" : ISODate("2018-10-31T00:00:00Z")}})
