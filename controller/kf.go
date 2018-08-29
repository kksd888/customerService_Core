// 客服相关

package controller

import (
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type KfServerController struct {
}

func InitKfServer() *KfServerController {
	return &KfServerController{}
}

// @Summary 获取客服信息
// @Description 获取客服信息
// @Tags Kf
// @Accept  json
// @Produce  json
// @Param id path int true "客服的ID"
// @Success 200 {string} json ""
// @Router /v1/server/{id} [get]
func (c *KfServerController) Get(context *gin.Context) {
}

// @Summary 客服修改在线状态
// @Description 客服修改在线状态
// @Tags Kf
// @Accept  json
// @Produce  json
// @Param id path int true "客服的ID"
// @Success 200 {string} json ""
// @Router /v1/server/{id}/status [put]
func (c *KfServerController) ChangeStatus(context *gin.Context) {
}

// @Summary 客服登入
// @Description 客服登入
// @Tags Kf
// @Accept  json
// @Produce  json
// @Param tokenId path int true "客服的授权TokenId"
// @Success 200 {string} json ""
// @Router /v1/login/:tokenId [post]
func (c *KfServerController) LoginIn(context *gin.Context) {
	tokenId := context.Param("tokenId")
	if tokenId == "" {
		context.JSON(http.StatusOK, gin.H{"code": 1000, "msg": "缺少授权客服的token"})
		return
	}

	kf := model.Kf{}
	if err := kf.GetByTokenId(tokenId); err != nil {
		log.Fatalf("LoginIn error：%#v", err)
	} else {
		// todo 向内存写入登录的客服
		s, _ := logic.RoomKf{
			KfId:         kf.Id,
			KfName:       kf.NickName,
			KfHeadImgUrl: kf.HeadImgUrl,
			KfStatus:     0,
		}.Make2Auth()
		context.JSON(http.StatusOK, LoginInResponse{
			BaseResponse: BaseResponse{
				Code: 0,
			},
			Authentication: s,
		})
	}
}

// 客服登出
func (c *KfServerController) LoginOut(context *gin.Context) {

}

type LoginInResponse struct {
	BaseResponse
	Authentication string
}
