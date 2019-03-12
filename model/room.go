package model

import (
	"customerService_Core/common"
	"time"
)

// 聊天室，在mongodb中用room文档来维护和存储对话数据，聊天数据也将保留100条，用来隔离用户，在并发下安全
type Room struct {
	RoomCustomer RoomCustomer  `json:"room_customer" bson:"room_customer"`
	RoomKf       RoomKf        `json:"room_kf" bson:"room_kf"`
	RoomMessages []RoomMessage `json:"room_messages" bson:"room_messages"`
	CreateTime   time.Time     `json:"create_time" bson:"create_time"`
}
type RoomCustomer struct {
	CustomerId           string                    `json:"customer_id" bson:"customer_id"`
	CustomerNickName     string                    `json:"customer_nick_name" bson:"customer_nick_name"`
	CustomerHeadImgUrl   string                    `json:"customer_head_img_url" bson:"customer_head_img_url"`
	CustomerSource       common.CustomerSourceType `json:"customer_source" bson:"customer_source"`
	CustomerPreviousKfId string                    `json:"customer_previous_kf_id" bson:"customer_previous_kf_id"`
}
type RoomKf struct {
	KfId         string `json:"kf_id" bson:"kf_id"`
	KfName       string `json:"kf_name" bson:"kf_name"`
	KfHeadImgUrl string `json:"kf_head_img_url" bson:"kf_head_img_url"`
}
type RoomMessage struct {
	Id         string    `json:"id" bson:"id"`
	Type       string    `json:"type" bson:"type"`
	MediaUrl   string    `json:"media_url" bson:"media_url"`
	Msg        string    `json:"msg" bson:"msg"`
	AiMsg      string    `json:"ai_msg" bson:"ai_msg"`
	Ack        bool      `json:"ack" bson:"ack"`
	OperCode   int       `json:"oper_code" bson:"oper_code"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}

func (r *Room) FormatterTimeLocation() {
	r.CreateTime = r.CreateTime.In(common.LocalLocation)
}

func (r *RoomMessage) FormatterTimeLocation() {
	r.CreateTime = r.CreateTime.In(common.LocalLocation)
}
