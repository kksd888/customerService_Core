package admin

import (
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"strconv"
	"time"
)

type StatisticsController struct {
}

func NewStatistics() *StatisticsController {
	return &StatisticsController{}
}

// @Summary 查询时间区间内接待的客户总量
// @Description 查询时间区间内接待的客户总量
// @Tags Dialog
// @Accept  json
// @Produce  json
// @Param kfid path int true "客服 ID"
// @page starTime path int true "开始时间"
// @limit endTime path int true "结束时间"
// @Success 200 {string} json ""
// @Router /admin/statistics/MessageCountByKf/{kfid}/{starTime}/{endTime} [get]
func (c *StatisticsController) MessageCountByKf(context *gin.Context) {
	var (
		kfid        = context.Param("kfid")
		starTimeStr = context.Param("starTime")
		endTimeStr  = context.Param("endTime")
	)
	if kfid == "" {
		//kfid = "06f17d3d66194b24a72a3400db3fb9e9"
		ReturnErrInfo(context, errors.New("缺少kfid"))
	}
	if starTimeStr == "" {
		//starTimeStr = "2018-01-01"
		ReturnErrInfo(context, errors.New("缺少开始时间"))
	}
	if endTimeStr == "" {
		//endTimeStr = "2018-12-31"
		ReturnErrInfo(context, errors.New("缺少结束时间"))
	}

	starTime, err := time.Parse("2006-01-02 15:04:05", starTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("开始时间格式错误 格式应为:2006-01-02 15:04:05"))
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("结束时间格式错误,格式应为:2006-01-02 15:04:05"))
	}

	var (
		query             = bson.M{"kf_id": kfid, "create_time": bson.M{"$gte": starTime, "$lt": endTime}}
		messageCollection = model.Db.C("message")
	)

	messageCount, _ := messageCollection.Find(query).Count()
	context.JSON(http.StatusOK, messageCount)
}

// @Router /admin/statistics/CustomerCountByKf/{kfid}/{starTime}/{endTime} [get]
func (c *StatisticsController) CustomerCountByKf(context *gin.Context) {
	var (
		kfid        = context.Param("kfid")
		starTimeStr = context.Param("starTime")
		endTimeStr  = context.Param("endTime")
	)
	if kfid == "" {
		kfid = "06f17d3d66194b24a72a3400db3fb9e9"
		ReturnErrInfo(context, errors.New("缺少kfid"))
	}
	if starTimeStr == "" {
		starTimeStr = "2018-09-07 00:00:00"
		ReturnErrInfo(context, errors.New("缺少开始时间"))
	}
	if endTimeStr == "" {
		endTimeStr = "2018-09-08 00:00:00"
		ReturnErrInfo(context, errors.New("缺少结束时间"))
	}

	starTime, err := time.Parse("2006-01-02 15:04:05", starTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("开始时间格式错误 格式应为:2006-01-02 15:04:05"))
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("结束时间格式错误,格式应为:2006-01-02 15:04:05"))
	}

	var (
		result []map[string]interface{}
		query  = []bson.M{
			{
				"$match": bson.M{"kf_id": kfid, "create_time": bson.M{"$gte": starTime, "$lt": endTime}},
			},
			{
				"$sort": bson.M{"customer_id": 1},
			},
			{
				"$skip": 1,
			},
			{
				"$limit": 100,
			},
			{
				"$group": bson.M{
					"_id":         "$_id",
					"customer_id": bson.M{"$first": "$customer_id"},
				},
			},
		}
		roomCollection = model.Db.C("message")
	)

	if err := roomCollection.Find(query).All(&result); err != nil {
		ReturnErrInfo(context, "未查到数据!")
	}
	context.JSON(http.StatusOK, len(result))
}

// @Router /admin/statistics/MessageCount/{kfid}/{starTime}/{endTime} [get]
func (c *StatisticsController) MessageCount(context *gin.Context) {
	var (
		kfid        = context.Param("kfid")
		starTimeStr = context.Param("starTime")
		endTimeStr  = context.Param("endTime")
		strPage     = context.Param("page")
		strLimit    = context.Param("limit")
	)
	if kfid == "" {
		//kfid = "06f17d3d66194b24a72a3400db3fb9e9"
		ReturnErrInfo(context, errors.New("缺少kfid"))
	}
	if starTimeStr == "" {
		//starTimeStr = "2018-01-01"
		ReturnErrInfo(context, errors.New("缺少开始时间"))
	}
	if endTimeStr == "" {
		//endTimeStr = "2018-12-31"
		ReturnErrInfo(context, errors.New("缺少结束时间"))
	}
	page, err := strconv.Atoi(strPage)
	if err != nil {
		ReturnErrInfo(context, errors.New("缺少page"))
	}
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		ReturnErrInfo(context, errors.New("缺少limit"))
	}

	starTime, err := time.Parse("2006-01-02 15:04:05", starTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("开始时间格式错误 格式应为:2006-01-02 15:04:05"))
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("结束时间格式错误,格式应为:2006-01-02 15:04:05"))
	}

	var (
		query = []bson.M{
			{
				"$match": bson.M{"kf_id": kfid, "create_time": bson.M{"$gte": starTime, "$lt": endTime}},
			},
			{"$lookup": bson.M{
				"from":         "kefu",
				"localField":   "kf_id",
				"foreignField": "id",
				"as":           "kf_info",
			}},
			{
				"$sort": bson.M{"create_time": 1},
			},
			{
				"$skip": (page - 1) * limit,
			},
			{
				"$limit": limit,
			},
		}
		messageCollection = model.Db.C("message")
	)

	messageCount, _ := messageCollection.Find(query).Count()
	context.JSON(http.StatusOK, messageCount)
}

func (c *StatisticsController) StatisticeMessage(context *gin.Context) {
	var (
		kfid        = context.Param("kfid")
		starTimeStr = context.Param("starTime")
		endTimeStr  = context.Param("endTime")
		strPage     = context.Param("page")
		strLimit    = context.Param("limit")
	)
	if kfid == "" {
		//kfid = "06f17d3d66194b24a72a3400db3fb9e9"
		ReturnErrInfo(context, errors.New("缺少kfid"))
	}
	if starTimeStr == "" {
		//starTimeStr = "2018-01-01"
		ReturnErrInfo(context, errors.New("缺少开始时间"))
	}
	if endTimeStr == "" {
		//endTimeStr = "2018-12-31"
		ReturnErrInfo(context, errors.New("缺少结束时间"))
	}
	page, err := strconv.Atoi(strPage)
	if err != nil {
		ReturnErrInfo(context, errors.New("缺少page"))
	}
	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		ReturnErrInfo(context, errors.New("缺少limit"))
	}

	starTime, err := time.Parse("2006-01-02 15:04:05", starTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("开始时间格式错误 格式应为:2006-01-02 15:04:05"))
	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {
		ReturnErrInfo(context, errors.New("结束时间格式错误,格式应为:2006-01-02 15:04:05"))
	}

	var (
		query = []bson.M{
			{
				"$match": bson.M{"create_time": bson.M{"$gte": starTime, "$lt": endTime}},
			},
			{"$lookup": bson.M{
				"from":         "kefu",
				"localField":   "kf_id",
				"foreignField": "id",
				"as":           "kf_info",
			}},
			{
				"$sort": bson.M{"create_time": 1},
			},
			{
				"$skip": (page - 1) * limit,
			},
			{
				"$limit": limit,
			},
		}
		messageCollection = model.Db.C("message")
	)

	messageCount, _ := messageCollection.Find(query).Count()
	context.JSON(http.StatusOK, messageCount)
}
