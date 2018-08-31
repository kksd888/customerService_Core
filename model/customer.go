package model

import (
	"time"
)

// 用户模型
type Customer struct {
	CustomerId   string
	NickName     string
	Sex          int32
	HeadImgUrl   string
	Address      string
	CustomerType int // 1、普通
	CreateTime   time.Time
	UpdateTime   time.Time
}
