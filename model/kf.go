package model

import (
	"time"
)

// 客服模型
type Kf struct {
	Id         string    `json:"id" bson:"id"`
	TokenId    string    `json:"token_id" bson:"token_id"`
	NickName   string    `json:"nick_name" bson:"nick_name"`
	Type       int       `json:"type" bson:"type"`
	HeadImgUrl string    `json:"head_img_url" bson:"head_img_url"`
	IsOnline   bool      `json:"is_online" bson:"is_online"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	UpdateTime time.Time `json:"update_time" bson:"update_time"`
}
