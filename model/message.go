package model

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/li-keli/go-tool/util/db_util"
	"log"
	"time"
)

// 历史存储的消息，存档，仅共查询
type Message struct {
	Id         string    `bson:"id"`
	CustomerId string    `bson:"customer_id"`
	KfId       string    `bson:"kf_id"`
	Type       string    `bson:"type"`
	MediaUrl   string    `bson:"media_url"`
	Msg        string    `bson:"msg"`
	AiMsg      string    `bson:"ai_msg"`
	OperCode   int       `bson:"oper_code"`
	CreateTime time.Time `bson:"create_time"`
}

func InsertMessage(m Message) {
	session := db_util.MongoDbSession.Copy()
	defer session.Close()

	if err := session.DB(common.AppConfig.DbName).C("message").Insert(&m); err != nil {
		log.Printf("消息存储异常：%s", err.Error())
	}
}
