package model

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"time"
)

// 用户模型
type Customer struct {
	CustomerId         string                    `bson:"customer_id" json:"customer_id"`
	NickName           string                    `bson:"nick_name" json:"nick_name"`
	Sex                int32                     `bson:"sex" json:"sex"`
	HeadImgUrl         string                    `bson:"head_img_url" json:"head_img_url"`
	Address            string                    `bson:"address" json:"address"`
	CustomerType       int                       `bson:"customer_type" json:"customer_type"`
	CustomerSourceType common.CustomerSourceType `json:"customer_source_type" bson:"customer_source_type"`
	CreateTime         time.Time                 `bson:"create_time" json:"create_time"`
	UpdateTime         time.Time                 `bson:"update_time" json:"update_time"`
}
