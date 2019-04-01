package open

import (
	"customerService_Core/common"
	"customerService_Core/handle"
	"customerService_Core/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/li-keli/go-tool/util/mongo_util"
	"github.com/li-keli/mgo/bson"
	"time"
)

type OpenController struct {
}

func NewOpen() *OpenController {
	return &OpenController{}
}

// 认证授权
// /v1/app/access
func (open *OpenController) Access(ctx *gin.Context) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		input = struct {
			DeviceId   string                    `json:"device_id" bson:"device_id" binding:"required"` // 设备编号
			CustomerId string                    `json:"customer_id" bson:"customer_id"`                // 用户编号
			NickName   string                    `json:"nick_name" bson:"nick_name"`                    // 用户昵称
			HeadImgUrl string                    `json:"head_img_url" bson:"head_img_url"`              // 用户头像
			Source     common.CustomerSourceType `json:"-" bson:"customer_source_type"`                 // 来源
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

		customerCollection = session.DB(common.AppConfig.DbName).C("customer")
		roomCollection     = session.DB(common.AppConfig.DbName).C("room")
		kefuCollection     = session.DB(common.AppConfig.DbName).C("kefu")
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
	if input.HeadImgUrl == "" {
		input.HeadImgUrl = common.RandomHeadImg()
	}

	lineMsg := ""
	onlineKefuCount, _ := kefuCollection.Find(bson.M{"is_online": true}).Count()
	lineCount, _ := roomCollection.Find(bson.M{"room_kf.kf_id": "", "room_messages.oper_code": common.MessageFromCustomer}).Count()
	if onlineKefuCount == 0 {
		lineMsg = common.KF_REPLY
	} else {
		if lineCount == 0 {
			lineMsg = common.WELCOME_REPLY
		} else {
			lineMsg = fmt.Sprintf(common.LINE_UP_REPLY, lineCount)
		}
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
					Msg:        lineMsg,
					OperCode:   common.MessageFromSys,
					Ack:        true,
					CreateTime: time.Now(),
				},
			},
			CreateTime: time.Now(),
		})
	} else {
		var (
			kefuModel = model.Kf{}
			room      = model.Room{}
		)
		roomCollection.Find(bson.M{"room_customer.customer_id": input.CustomerId}).One(&room)
		kefuCollection.Find(bson.M{"id": room.RoomKf.KfId}).One(&kefuModel)
		if kefuModel.Id != "" && kefuModel.IsOnline == false {
			// 若接待的客服已经下线，则将用户重新放入待接入
			_ = roomCollection.Update(
				bson.M{"room_customer.customer_id": input.CustomerId},
				bson.M{"$set": bson.M{"room_kf": &model.RoomKf{}}})

			// 更新默认欢迎消息
			query := bson.M{
				"room_customer.customer_id": input.CustomerId,
			}
			changes := bson.M{
				"$push": bson.M{"room_messages": bson.M{"$each": []model.RoomMessage{
					{
						Id:         common.GetNewUUID(),
						Type:       string(common.MsgTypeText),
						Msg:        lineMsg,
						OperCode:   common.MessageFromSys,
						Ack:        true,
						CreateTime: time.Now(),
					},
				},
					"$slice": -100}},
			}
			_ = roomCollection.Update(query, changes)
		}
		if kefuModel.Id == "" {
			// 更新默认欢迎消息
			query := bson.M{
				"room_customer.customer_id": input.CustomerId,
			}
			changes := bson.M{
				"$push": bson.M{"room_messages": bson.M{"$each": []model.RoomMessage{
					{
						Id:         common.GetNewUUID(),
						Type:       string(common.MsgTypeText),
						Msg:        lineMsg,
						OperCode:   common.MessageFromSys,
						Ack:        true,
						CreateTime: time.Now(),
					},
				},
					"$slice": -100}},
			}
			_ = roomCollection.Update(query, changes)
		}
	}

	// 生成授权码
	auth, err := handle.OpenMake2Auth(input.CustomerId)
	if err != nil {
		common.ReturnErr(ctx, err)
	}
	output.Authorization = auth

	common.ReturnSuccess(ctx, output)
}
