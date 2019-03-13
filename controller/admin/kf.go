// 客服相关

package admin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/li-keli/go-tool/util/db_util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	session := db_util.MongoDbSession.Copy()
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
	session := db_util.MongoDbSession.Copy()
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
	session := db_util.MongoDbSession.Copy()
	defer session.Close()
	var (
		kf           = model.Kf{}
		output       = LoginEmployeeResponse{}
		openId, _    = context.Get("OpenID")
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

	//请求会员登录接口
	output = GetEmployeeInfo(loginStruct.JobNum, loginStruct.PassWord, openId)
	if output.BaseResponse.IsSuccess {
		if err := kfCollection.Find(bson.M{
			"job_num": loginStruct.JobNum,
		}).One(&kf); err != nil {
			//添加用户
			kf.Id = common.ToMd5(loginStruct.JobNum + loginStruct.GroupName)
			kfCollection.Insert(&model.Kf{
				Id:         kf.Id,
				JobNum:     loginStruct.JobNum,
				NickName:   loginStruct.JobNum,
				IsOnline:   true,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
				Type:       1,
				GroupName:  loginStruct.GroupName,
			})
		} else {
			// 更新在线客服列表
			changeKfModel := model.Kf{Id: kf.Id, IsOnline: true}
			err := changeKfModel.ChangeStatus()
			if err != nil {
				ReturnErrInfo(context, err)
			}
		}
	} else {
		ReturnErrInfo(context, "用户名呼和密码错误")
	}

	s, _ := Make2Auth(kf.Id)

	context.JSON(http.StatusOK, LoginInResponse{
		Authentication: s,
	})
}

func Make2Auth(kfId string) (string, error) {
	encrypt := common.AesEncrypt{}
	byteInfo, err := encrypt.Encrypt([]byte(kfId))
	if err != nil {
		log.Printf("common.NewGoAES() err：%v", err)
	}

	return base64.StdEncoding.EncodeToString(byteInfo), err
}

func GetEmployeeInfo(employeeAlias string, password string, openId string) LoginEmployeeResponse {
	var (
		output = LoginEmployeeResponse{}
	)

	var input = struct {
		MethodName    string `json:"MethodName"`
		EmployeeAlias string `json:"EmployeeAlias"`
		Password      string `json:"Password"`
		LoginDeviceID string `json:"LoginDeviceID"`
	}{
		MethodName:    LoginEmployeeMonth,
		EmployeeAlias: employeeAlias,
		Password:      password,
		LoginDeviceID: openId,
	}
	marshal, _ := json.Marshal(input)

	req, err := http.NewRequest("POST", "http://memberapi.jsjinfo.cn/Hosts/JIUser.aspx", bytes.NewBuffer(marshal))
	req.Header.Set("MethodName", LoginEmployeeMonth)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &output); err != nil {
		logrus.Error(err)
	}

	if !output.BaseResponse.IsSuccess {
		logrus.Error(output.BaseResponse.ErrorMessage)
	}

	return output
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
	Authentication string
}
