package model

import "time"

type Room struct {
	RoomCustomer RoomCustomer  `json:"room_customer" bson:"room_customer"`
	RoomKf       RoomKf        `json:"room_kf" bson:"room_kf"`
	RoomMessages []RoomMessage `json:"room_messages" bson:"room_messages"`
	CreateTime   time.Time     `json:"create_time" bson:"create_time"`
}
type RoomCustomer struct {
	CustomerId           string `json:"customer_id" bson:"customer_id"`
	CustomerNickName     string `json:"customer_nick_name" bson:"customer_nick_name"`
	CustomerHeadImgUrl   string `json:"customer_head_img_url" bson:"customer_head_img_url"`
	CustomerPreviousKfId string `json:"customer_previous_kf_id" bson:"customer_previous_kf_id"`
}
type RoomKf struct {
	KfId         string `json:"kf_id" bson:"kf_id"`
	KfName       string `json:"kf_name" bson:"kf_name"`
	KfHeadImgUrl string `json:"kf_head_img_url" bson:"kf_head_img_url"`
	KfStatus     int    `json:"kf_status" bson:"kf_status"`
}
type RoomMessage struct {
	Id         string    `json:"id" bson:"id"`
	Type       string    `json:"type" bson:"type"`
	MediaUrl   string    `json:"media_url" bson:"media_url"`
	Msg        string    `json:"msg" bson:"msg"`
	Ack        bool      `json:"ack" bson:"ack"`
	OperCode   int       `json:"oper_code" bson:"oper_code"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
}
