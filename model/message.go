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
	MsgType       int       `gorm:"msg_type"`
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
	return "chat_message"
}

func (m Message) InsertOrUpdate() {
	m.CreateTime = time.Now()
	m.UpdateTime = time.Now()
	createErr := db.Create(&m)
	if createErr.Error != nil {
		panic(createErr.Error)
	}
}

func (m Message) AccessAck() {
	update := db.Model(&m).Where("kf_id=? and customer_token=?", m.KfId, m.CustomerToken).Update("kf_ack")
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
	find := db.Raw(`select chat_message.id,
       chat_message.customer_token,
       chat_message.kf_id,
       chat_message.msg,
       chat_message.msg_type,
       chat_message.oper_code,
       chat_message.kf_ack,
       chat_message.create_time,
       chat_message.update_time,
       dic_customer.nick_name as customer_nick_name,
       dic_customer.head_img_url as customer_head_img_url
from chat_message
       left join dic_customer on dic_customer.open_id = chat_message.customer_token
where kf_id = ? limit 10;`, m.KfId).Scan(&messages)
	return messages, find.Error
}
