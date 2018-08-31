package model

import (
	"time"
)

// 消息模型
type Message struct {
	Id            int       `gorm:"id"`
	CustomerToken string    `gorm:"customer_token"`
	KfId          int       `gorm:"kf_id"`
	KfAck         bool      `gorm:"kf_ack"`
	Msg           string    `gorm:"msg"`
	MsgType       string    `gorm:"msg_type"`
	OperCode      int       `gorm:"oper_code"`
	CreateTime    time.Time `gorm:"create_time"`
	UpdateTime    time.Time `gorm:"update_time"`
}

// 消息模型-扩展用户数据
type MessageLinkCustomer struct {
	Message
	CustomerNickName   string `gorm:"customer_nick_name"`
	CustomerHeadImgUrl string `gorm:"customer_head_img_url"`
}
