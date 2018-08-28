package model

import "time"

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

func (m Message) InsertOrUpdate() error {
	create := db.Create(&m)
	return create.Error
}
