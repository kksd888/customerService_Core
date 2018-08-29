package model

import (
	"log"
	"time"
)

// 消息模型
type Message struct {
	Id         int       `gorm:"id"`
	CustomerId int       `gorm:"customer_id"`
	KfId       int       `gorm:"kf_id"`
	Msg        string    `gorm:"msg"`
	MsgType    int       `gorm:"msg_type"`
	OperCode   int       `gorm:"oper_code"`
	CreateTime time.Time `gorm:"create_time"`
	UpdateTime time.Time `gorm:"update_time"`
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

func (m *Message) InsertOrUpdate() error {
	create := db.Create(m)
	return create.Error
}

func (m *MessageLinkCustomer) GetKfHistoryMsg() ([]MessageLinkCustomer, error) {
	var messages []MessageLinkCustomer
	find := db.Raw(`select chat_message.id,
       chat_message.customer_id,
       chat_message.kf_id,
       chat_message.msg,
       chat_message.msg_type,
       chat_message.oper_code,
       chat_message.create_time,
       chat_message.update_time,
       dic_customer.nick_name as customer_nick_name,
       dic_customer.head_img_url as customer_head_img_url
from chat_message
       left join dic_customer on dic_customer.id = chat_message.customer_id
where kf_id = ? limit 10;`, m.KfId).Scan(&messages)
	return messages, find.Error
}
