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

func (Message) TableName() string {
	return "dialog_message"
}

func (m Message) Insert() {
	m.CreateTime = time.Now()
	m.UpdateTime = time.Now()
	createErr := db.Create(&m)
	if createErr.Error != nil {
		panic(createErr.Error)
	}
}

func (m Message) Access() {
	update := db.Model(&m).Where("kf_id=0 and customer_token=?", m.CustomerToken).Update("kf_id", m.KfId)
	if update.Error != nil {
		panic(update.Error)
	}
}

func (m Message) Ack() {
	update := db.Model(&m).Where("kf_id=? and customer_token=?", m.KfId, m.CustomerToken).Update("kf_ack", m.KfAck)
	if update.Error != nil {
		panic(update.Error)
	}
}

func (m *Message) WaitReply() ([]Message, error) {
	var messages []Message
	find := db.Where("kf_id=? and kf_ack=?", m.KfId, false).Find(&messages)
	return messages, find.Error
}

func (m *MessageLinkCustomer) GetKfHistoryMsg() ([]MessageLinkCustomer, error) {
	var messages []MessageLinkCustomer
	find := db.Raw(`select dialog_message.id,
       dialog_message.customer_token,
       dialog_message.kf_id,
       dialog_message.msg,
       dialog_message.msg_type,
       dialog_message.oper_code,
       dialog_message.kf_ack,
       dialog_message.create_time,
       dialog_message.update_time,
       dic_customer.nick_name as customer_nick_name,
       dic_customer.head_img_url as customer_head_img_url
from dialog_message
       left join dic_customer on dic_customer.open_id = dialog_message.customer_token
where kf_id = ? limit 10;`, m.KfId).Scan(&messages)
	return messages, find.Error
}
