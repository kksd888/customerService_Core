package model

import (
	"time"
)

// 客服模型
type Kf struct {
	Id         string    `bson:"id"`
	TokenId    string    `bson:"token_id"`
	NickName   string    `bson:"nick_name"`
	Type       int       `bson:"type"`
	HeadImgUrl string    `bson:"head_img_url"`
	Status     bool      `bson:"status"`
	CreateTime time.Time `bson:"create_time"`
	UpdateTime time.Time `bson:"update_time"`
}
