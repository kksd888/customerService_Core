// 客服相关

package controller

import (
	"encoding/base64"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
)

type KfServerController struct {
	db *model.MongoDb
}

func InitKfServer(_db *model.MongoDb) *KfServerController {
	return &KfServerController{db: _db}
}

// @Summary 获取客服信息
// @Description 获取客服信息
// @Tags Kf
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /v1/kf/{kfId} [get]
func (c *KfServerController) Get(context *gin.Context) {
	var (
		kf      model.Kf
		kfId, _ = context.Get("KFID")
		kfC     = c.db.C("kf")
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
// @Router /v1/kf/status [post]
func (c *KfServerController) ChangeStatus(context *gin.Context) {
	var (
		kfId, _ = context.Get("KFID")
		kfC     = c.db.C("kf")
		reqBind = struct {
			status bool `json:"status"`
		}{}
	)

	if err := context.Bind(reqBind); err != nil {
		ReturnErrInfo(context, err)
	}

	if err := kfC.Update(bson.M{"id": kfId}, bson.M{"status": reqBind.status}); err != nil {
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
// @Param tokenId path int true "客服的授权TokenId"
// @Success 200 {string} json "{"code":0,"msg":"ok"}"
// @Router /v1/login/{tokenId} [post]
func (c *KfServerController) LoginIn(context *gin.Context) {
	var (
		kf           = model.Kf{}
		tokenId      = context.Param("tokenId")
		kfCollection = c.db.C("kf")
	)

	if tokenId == "" {
		context.JSON(http.StatusOK, gin.H{"code": http.StatusUnauthorized, "msg": "缺少授权客服的token"})
		return
	}
	if err := kfCollection.Find(bson.M{"token_id": tokenId}).One(&kf); err != nil {
		ReturnErrInfo(context, err)
	}

	if kf.Id == "0" {
		ReturnErrInfo(context, errors.New("客服登录授权失败"))
	}

	//logic.AddOnlineKf(kf)
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

type LoginInResponse struct {
	Authentication string
}
