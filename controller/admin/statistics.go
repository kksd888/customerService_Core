package admin

import (
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/li-keli/go-tool/util/db_util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
)

type StatisticsController struct {
}

func NewStatistics() *StatisticsController {
	return &StatisticsController{}
}

// 统计查询
func (c *StatisticsController) Statistics(context *gin.Context) {
	session := db_util.MongoDbSession.Copy()
	defer session.Close()

	var (
		starTimeStr = context.Param("starTime")
		endTimeStr  = context.Param("endTime")
		pageStr     = context.Param("page")
		limitStr    = context.Param("limit")
	)
	starTime, err := time.Parse("2006-01-02 15:04:05", starTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("开始时间格式错误 格式应为 yyyy-MM-dd HH:mm:ss 如:2006-01-02 15:04:05"))
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("结束时间格式错误,格式应为 yyyy-MM-dd HH:mm:ss 如:2006-01-02 15:04:05"))
	}
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("page类型错误,应为数字,如 1"))
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("limit类型错误,应为数字,如 1000"))
	}

	//声明mongodb查询
	var (
		queryMessage = []bson.M{
			{
				"$match": bson.M{"kf_id": bson.M{"$ne": ""}, "create_time": bson.M{"$gte": starTime, "$lt": endTime}},
			},
			{"$lookup": bson.M{
				"from":         "kefu",
				"localField":   "kf_id",
				"foreignField": "id",
				"as":           "kefu",
			}},
			{
				"$unwind": bson.M{
					"path":                       "$kefu",
					"preserveNullAndEmptyArrays": true,
				},
			},
			{
				"$sort": bson.M{"kf_id": 1},
			},
			{
				"$group": bson.M{
					"_id":          "$kf_id",
					"kfId":         bson.M{"$first": "$kf_id"},
					"fkName":       bson.M{"$first": "$kefu.nick_name"},
					"messageCount": bson.M{"$sum": 1},
				},
			},
			{
				"$skip": (page - 1) * limit,
			},
			{
				"$limit": limit,
			},
		}
		queryCustomer = []bson.M{
			{
				"$match": bson.M{"kf_id": bson.M{"$ne": ""}, "create_time": bson.M{"$gte": starTime, "$lt": endTime}},
			},
			{
				"$sort": bson.M{"kf_id": 1},
			},
			{
				"$group": bson.M{
					"_id":           bson.M{"kf_id": "$kf_id", "customer_id": "$customer_id"},
					"kfId":          bson.M{"$first": "$kf_id"},
					"customerId":    bson.M{"$first": "$customer_id"},
					"customerCount": bson.M{"$sum": 1},
				},
			},
		}
		messageCollection = session.DB(common.DB_NAME).C("message")
	)

	//查询每个客服回复的信息
	var messageByKf []map[string]interface{}
	if err := messageCollection.Pipe(queryMessage).All(&messageByKf); err != nil {
		log.Warn(err)
	}

	//查询每个客服回复的用户
	var customerByKf []bson.M
	if err := messageCollection.Pipe(queryCustomer).All(&customerByKf); err != nil {
		log.Warn(err)
	}

	//循环赋值 每个客服回复的客服数量
	count := 0
	for i := 0; i < len(messageByKf); i++ {
		kfId := messageByKf[i]["kfId"].(string)
		count = 0
		for j := 0; j < len(customerByKf); j++ {
			if kfId == customerByKf[j]["kfId"].(string) {
				count++
			}
		}
		messageByKf[i]["customerCount"] = count
	}
	context.JSON(http.StatusOK, messageByKf)
}
