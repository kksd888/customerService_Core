// 客服相关

package admin

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
	var (
		kf      model.Kf
		kfId, _ = context.Get("KFID")
		kfC     = model.Db.C("kefu")
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
	var (
		kfId, _ = context.Get("KFID")
		kfC     = model.Db.C("kefu")
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

	output := GetEmployeeInfo("6094", "awk4js")
	ReturnErrInfo(context, output)

	//var (
	//	kf           = model.Kf{}
	//	kfCollection = model.Db.C("kefu")
	//	loginStruct  = struct {
	//		JobNum   string `json:"job_num"`
	//		PassWord string `json:"pass_word"`
	//	}{}
	//)
	//
	//if err := context.Bind(&loginStruct); err != nil {
	//	ReturnErrInfo(context, errors.New(fmt.Sprintf("登录参数错误：%s", err.Error())))
	//}
	//
	//if loginStruct.JobNum == "" || loginStruct.PassWord == "" {
	//	ReturnErrInfo(context, "登录参数错误")
	//}
	//
	//if err := kfCollection.Find(bson.M{
	//	"job_num":   loginStruct.JobNum,
	//	"pass_word": common.ToMd5(loginStruct.PassWord),
	//}).One(&kf); err != nil {
	//	ReturnErrInfo(context, "客服登录授权失败")
	//}
	//
	//s, _ := Make2Auth(kf.Id)
	//
	//// 更新在线客服列表
	//changeKfModel := model.Kf{Id: kf.Id, IsOnline: true}
	//err := changeKfModel.ChangeStatus()
	//if err != nil {
	//	ReturnErrInfo(context, err)
	//}
	//
	//context.JSON(http.StatusOK, LoginInResponse{
	//	Authentication: s,
	//})
}

func Make2Auth(kfId string) (string, error) {
	encrypt := common.AesEncrypt{}
	byteInfo, err := encrypt.Encrypt([]byte(kfId))
	if err != nil {
		log.Printf("common.NewGoAES() err：%v", err)
	}

	return base64.StdEncoding.EncodeToString(byteInfo), err
}

func GetEmployeeInfo(employeeAlias string, password string) LoginEmployeeResponse {
	var (
		output = LoginEmployeeResponse{}
	)

	// 计算签名
	//jsonBodyStr, _ := json.Marshal(jsonBody)
	//signStr := fmt.Sprintf("%s%s%d%s", strings.Replace(string(jsonBodyStr), " ", "", -1), LoginEmployeeMonth, nowUnixTime, key)
	//sign := util.ToMd5(signStr)

	var input = struct {
		MethodName    string `json:"MethodName"`
		EmployeeAlias string `json:"EmployeeAlias"`
		Password      string `json:"Password"`
		LoginDeviceID string `json:"LoginDeviceID"`
	}{
		MethodName:    LoginEmployeeMonth,
		EmployeeAlias: employeeAlias,
		Password:      password,
		LoginDeviceID: "862258037809663",
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
	IsVIPManager     string `json:"IsVIPManager"`
}

type LoginInResponse struct {
	Authentication string
}
