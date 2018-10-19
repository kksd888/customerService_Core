package open

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"net/http"
	"time"
)

type OpenController struct {
}

func NewOpen() *OpenController {
	return &OpenController{}
}

// 健康检查
// /health
func (c *OpenController) Health(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{"code": 100})
}

// 认证授权
// /v1/app/access
func (open *OpenController) Access(ctx *gin.Context) {
	var (
		input = struct {
			DeviceId   string                    `json:"device_id" bson:"device_id" binding:"required"` // 设备编号
			CustomerId string                    `json:"customer_id" bson:"customer_id"`                // 用户编号
			NickName   string                    `json:"nick_name" bson:"nick_name"`                    // 用户昵称
			HeadImgUrl string                    `json:"head_img_url" bson:"head_img_url"`              // 用户头像
			Source     common.CustomerSourceType `json:"-" bson:"source"`                               // 来源
			CreateTime time.Time                 `json:"-" bson:"create_time"`                          // DB创建时间
			UpdateTime time.Time                 `json:"-" bson:"update_time"`                          // DB更新时间
		}{
			Source:     common.FromAPP,
			CreateTime: time.Now(),
			UpdateTime: time.Now(),
		}
		output struct {
			Authorization string `json:"authorization"` // 授权码
		}

		customerCollection = model.Db.C("customer")
		roomCollection     = model.Db.C("room")
	)

	// 验证并绑定到模型
	if err := ctx.BindJSON(&input); err != nil {
		common.ReturnErrCode(ctx, common.ParameterBad, err)
	}

	// 游客配置
	if input.CustomerId == "" {
		input.CustomerId = input.DeviceId
	}
	if input.NickName == "" {
		input.NickName = "游客"
	}

	// 存储用户信息
	changeInfo, _ := customerCollection.Upsert(bson.M{"customer_id": input.CustomerId}, input)
	if changeInfo.Matched == 0 {
		// 更新默认欢迎消息
		roomCollection.Insert(&model.Room{
			RoomCustomer: model.RoomCustomer{
				CustomerId:         input.CustomerId,
				CustomerNickName:   input.NickName,
				CustomerHeadImgUrl: input.HeadImgUrl,
				CustomerSource:     input.Source,
			},
			RoomMessages: []model.RoomMessage{
				{
					Id:         common.GetNewUUID(),
					Type:       string(common.MsgTypeText),
					Msg:        "你好",
					OperCode:   common.MessageFromCustomer,
					CreateTime: time.Now(),
				},
			},
			CreateTime: time.Now(),
		})
	}

	// 生成授权码
	auth, err := handle.OpenMake2Auth(input.CustomerId)
	if err != nil {
		common.ReturnErr(ctx, err)
	}
	output.Authorization = auth

	common.ReturnSuccess(ctx, output)
}
