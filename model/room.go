package model

import "time"

type Room struct {
	RoomCustomer RoomCustomer  `bson:"room_customer"`
	RoomKf       RoomKf        `bson:"room_kf"`
	RoomMessages []RoomMessage `bson:"room_messages"`
	CreateTime   time.Time     `bson:"create_time"`
}
type RoomCustomer struct {
	CustomerId           string `bson:"customer_id"`
	CustomerNickName     string `bson:"customer_nick_name"`
	CustomerHeadImgUrl   string `bson:"customer_head_img_url"`
	CustomerPreviousKfId string `bson:"customer_previous_kf_id"`
}
type RoomKf struct {
	KfId         string `bson:"kf_id"`
	KfName       string `bson:"kf_name"`
	KfHeadImgUrl string `bson:"kf_head_img_url"`
	KfStatus     int    `bson:"kf_status"`
}
type RoomMessage struct {
	Id         string    `bson:"id"`
	Type       string    `bson:"type"`
	Msg        string    `bson:"msg"`
	Ack        bool      `bson:"ack"`
	OperCode   int       `bson:"oper_code"`
	CreateTime time.Time `bson:"create_time"`
}
