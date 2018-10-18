// 离线自动回复设置

package admin

import (
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
)

type OfflineReplyController struct {
}

func NewOfflineReply() *OfflineReplyController {
	return &OfflineReplyController{}
}

// @Summary 获取自动回复列表
// @Description 获取自动回复列表
// @Tags OfflineReply
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /admin/setting/offline_reply [get]
func (c *OfflineReplyController) List(context *gin.Context) {
}

// @Summary 新增一条离线自动回复语句
// @Description 新增一条离线自动回复语句
// @Tags OfflineReply
// @Accept  json
// @Produce  json
// @Success 200 {string} json ""
// @Router /admin/setting/offline_reply [post]
func (c *OfflineReplyController) Create(context *gin.Context) {
	var (
		offline = struct {
			offLineMsg string `bson:"off_line_msg" json:"off_line_msg"`
			operKfId   string `bson:"oper_kf_id" json:"oper_kf_id"`
		}{}
		replyC = model.Db.C("offineReply")
	)

	if err := replyC.Insert(&offline); err != nil {
		ReturnErrInfo(context, err)
	}
}

// @Summary 删除一条离线自动回复语句
// @Description 删除一条离线自动回复语句
// @Tags OfflineReply
// @Accept  json
// @Produce  json
// @Param id path int true "自动回复语句的ID"
// @Success 200 {string} json ""
// @Router /admin/setting/offline_reply/{id} [delete]
func (c *OfflineReplyController) Delete(context *gin.Context) {
}

// @Summary 更新一条离线自动回复语句
// @Description 更新一条离线自动回复语句
// @Tags OfflineReply
// @Accept  json
// @Produce  json
// @Param id path int true "自动回复语句的ID"
// @Success 200 {string} json ""
// @Router /admin/setting/offline_reply/{id} [put]
func (c *OfflineReplyController) Update(context *gin.Context) {
}
