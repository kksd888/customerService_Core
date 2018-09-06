package model

import (
	"log"
	"time"
)

// 历史存储的消息
type Message struct {
	Id         string    `bson:"id"`
	CustomerId string    `bson:"customer_id"`
	KfId       string    `bson:"kf_id"`
	Type       string    `bson:"type"`
	MediaUrl   string    `bson:"media_url"`
	Msg        string    `bson:"msg"`
	OperCode   int       `bson:"oper_code"`
	CreateTime time.Time `bson:"create_time"`
}

func InsertMessage(m Message) {
	if err := Db.C("message").Insert(&m); err != nil {
		log.Printf("消息存储异常：%s", err.Error())
	}
}
