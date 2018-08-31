package model

import (
	"time"
)

// 客服模型
type Kf struct {
	Id         string
	TokenId    string
	NickName   string
	Type       int
	HeadImgUrl string
	CreateTime time.Time
	UpdateTime time.Time
}
