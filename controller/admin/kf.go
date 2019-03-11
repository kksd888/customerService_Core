// 客服相关

package admin

import (
	"encoding/base64"
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type KfServerController struct {
}

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
	session := model.DbSession.Copy()
	defer session.Close()

	var (
		kf      model.Kf
		kfId, _ = context.Get("KFID")
		kfC     = session.DB(common.DB_NAME).C("kefu")
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
	session := model.DbSession.Copy()
	defer session.Close()

	var (
		kfId, _ = context.Get("KFID")
		kfC     = session.DB(common.DB_NAME).C("kefu")
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
	session := model.DbSession.Copy()
	defer session.Close()

	var (
		kf           = model.Kf{}
		kfCollection = session.DB(common.DB_NAME).C("kefu")
		loginStruct  = struct {
			JobNum   string `json:"job_num"`
			PassWord string `json:"pass_word"`
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
	}

	s, _ := Make2Auth(kf.Id)

	// 更新在线客服列表
	changeKfModel := model.Kf{Id: kf.Id, IsOnline: true}
	err := changeKfModel.ChangeStatus()
	if err != nil {
		ReturnErrInfo(context, err)
	}

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

type LoginInResponse struct {
	Authentication string
}
