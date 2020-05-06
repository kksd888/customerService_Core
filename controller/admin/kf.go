// 客服相关

package admin

import (
	"customerService_Core/common"
	"customerService_Core/model"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/li-keli/go-tool/util/mongo_util"
	"github.com/li-keli/mgo/bson"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type KfServerController struct {
}

var (
	LoginEmployeeMonth = "LoginEmployee"
)

func NewKfServer() *KfServerController {
	return &KfServerController{}
}

// @Summary 获取客服信息
// @Description 获取客服信息
// @Tags Kf
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /admin/kf [get]
func (c *KfServerController) Get(context *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		kf      model.Kf
		kfId, _ = context.Get("KFID")
		kfC     = session.DB(common.AppConfig.DbName).C("kefu")
	)

	if err := kfC.Find(bson.M{"id": kfId}).One(&kf); err != nil {
		ReturnErrInfo(context, err)
	}

	context.JSON(http.StatusOK, kf)
}

// @Summary 客服修改在线状态
// @Description 客服修改在线状态
// @Tags Kf
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /admin/kf/status [post]
func (c *KfServerController) ChangeStatus(context *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		kfId, _ = context.Get("KFID")
		kfC     = session.DB(common.AppConfig.DbName).C("kefu")
		reqBind = struct {
			Status bool `bson:"status" json:"status"`
		}{}
	)

	if err := context.Bind(&reqBind); err != nil {
		ReturnErrInfo(context, err)
	}

	if err := kfC.Update(bson.M{"id": kfId}, bson.M{"$set": bson.M{"is_online": reqBind.Status}}); err != nil {
		ReturnErrInfo(context, err)
	} else {
		ReturnSuccessInfo(context)
	}
}

// @Summary 客服登入
// @Description 客服登入
// @Tags Kf
// @Accept  json
// @Produce  json
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /admin/login [post]
func (c *KfServerController) LoginIn(context *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()
	var (
		kf = model.Kf{}
		//memberOutApi = LoginEmployeeResponse{}
		kfCollection = session.DB(common.AppConfig.DbName).C("kefu")
		loginStruct  = struct {
			JobNum    string `json:"job_num"`
			PassWord  string `json:"pass_word"`
			GroupName string `json:"group_name"`
		}{}
	)

	if err := context.Bind(&loginStruct); err != nil {
		ReturnErrInfo(context, errors.New(fmt.Sprintf("登录参数错误：%s", err.Error())))
	}

	if loginStruct.JobNum == "" || loginStruct.PassWord == "" {
		ReturnErrInfo(context, "登录参数错误")
	}

	if err := kfCollection.Find(bson.M{
		"job_num":   loginStruct.JobNum,
		"pass_word": common.ToMd5(loginStruct.PassWord),
	}).One(&kf); err != nil {
		ReturnErrInfo(context, "客服登录授权失败")
	} else {
		// 更新在线客服列表
		if err := kfCollection.Update(bson.M{"job_num": loginStruct.JobNum},
			bson.M{"$set": bson.M{
				"is_online":  true,
				"nick_name":  kf.NickName,
				"group_name": loginStruct.GroupName,
			}}); err != nil {
			ReturnErrInfo(context, err)
		}
	}

	////请求会员登录接口
	//if memberOutApi = GetEmployeeInfo(loginStruct.JobNum, loginStruct.PassWord, "wechar_kf"); memberOutApi.BaseResponse.IsSuccess {
	//	if err := kfCollection.Find(bson.M{"job_num": loginStruct.JobNum}).One(&kf); err != nil {
	//		//添加用户
	//		kf.Id = common.ToMd5(loginStruct.JobNum + loginStruct.GroupName)
	//		kfCollection.Insert(&model.Kf{
	//			Id:         kf.Id,
	//			JobNum:     loginStruct.JobNum,
	//			NickName:   memberOutApi.EmployeeName,
	//			IsOnline:   true,
	//			Type:       1,
	//			GroupName:  loginStruct.GroupName,
	//			CreateTime: time.Now(),
	//			UpdateTime: time.Now(),
	//		})
	//	} else {
	//		// 更新在线客服列表
	//		if err := kfCollection.Update(bson.M{"job_num": loginStruct.JobNum},
	//			bson.M{"$set": bson.M{
	//				"is_online":  true,
	//				"nick_name":  memberOutApi.EmployeeName,
	//				"group_name": loginStruct.GroupName,
	//			}}); err != nil {
	//			ReturnErrInfo(context, err)
	//		}
	//	}
	//} else {
	//	ReturnErrInfo(context, "用户名或密码错误")
	//}

	s, _ := Make2Auth(kf.Id)

	context.JSON(http.StatusOK, LoginInResponse{
		Authentication: s,
		NickName:       kf.NickName,
		JobNum:         loginStruct.JobNum,
		GroupName:      loginStruct.GroupName,
	})
}

// 在线客服列表
func (c *KfServerController) OnLines(ctx *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		kfModels     []bson.M
		kfId         = ctx.GetString("KFID")
		kfCollection = session.DB(common.AppConfig.DbName).C("kefu")
	)

	query := []bson.M{
		{
			"$match": bson.M{"is_online": true, "id": bson.M{"$ne": kfId}},
		},
		{
			"$group": bson.M{
				"_id":     "$group_name",
				"label":   bson.M{"$first": "$group_name"},
				"options": bson.M{"$push": bson.M{"value": "$id", "label": "$nick_name"}},
			},
		},
		{
			"$project": bson.M{
				"_id": 0,
			},
		},
	}
	_ = kfCollection.Pipe(query).All(&kfModels)

	ctx.JSON(http.StatusOK, kfModels)
}

func Make2Auth(kfId string) (string, error) {
	encrypt := common.AesEncrypt{}
	byteInfo, err := encrypt.Encrypt([]byte(kfId))
	if err != nil {
		log.Printf("common.NewGoAES() err：%v", err)
	}

	return base64.StdEncoding.EncodeToString(byteInfo), err
}

type LoginEmployeeResponse struct {
	BaseResponse struct {
		IsSuccess    bool   `json:"IsSuccess"`
		ErrorMessage string `json:"ErrorMessage"`
		ErrorCode    string `json:"ErrorCode"`
	}
	EmployeeName     string `json:"EmployeeName"`
	PositionID       int    `json:"PositionID"`
	DeptID           int    `json:"DeptID"`
	EmployeeID       int    `json:"EmployeeID"`
	IDNumber         string `json:"IDNumber"`
	Birthday         string `json:"Birthday"`
	NativePlace      string `json:"NativePlace"`
	Sex              int    `json:"Sex"`
	MobilePhone      string `json:"MobilePhone"`
	CompanyTel       string `json:"CompanyTel"`
	Memo             string `json:"Memo"`
	IsAdmin          int    `json:"IsAdmin"`
	PhoneSkill       string `json:"PhoneSkill"`
	DataCommission   int    `json:"DataCommission"`
	CallCenterNumber int    `json:"CallCenterNumber"`
	LastDeviceID     string `json:"LastDeviceID"`
	AvatarUrl        string `json:"AvatarUrl"`
	CompanyID        int    `json:"CompanyID"`
	Token            string `json:"token"`
	DeptName         string `json:"DeptName"`
	IsVIPManager     bool   `json:"IsVIPManager"`
}

type LoginInResponse struct {
	Authentication string `json:"Authentication"`
	JobNum         string `json:"JobNum"`
	NickName       string `json:"NickName"`
	GroupName      string `json:"GroupName"`
}
