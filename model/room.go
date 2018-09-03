package model

import "time"

type Room struct {
	RoomCustomer
	RoomKf
	RoomMessages []RoomMessage
	CreateTime   time.Time
}
type RoomCustomer struct {
	CustomerId           string
	CustomerNickName     string
	CustomerHeadImgUrl   string
	CustomerPreviousKfId string
}
type RoomKf struct {
	KfId         string
	KfName       string
	KfHeadImgUrl string
	KfStatus     int
}
type RoomMessage struct {
	Id         string    `bson:"id"`
	Type       string    `bson:"type"`
	Msg        string    `bson:"msg"`
	Ack        bool      `bson:"ack"`
	OperCode   int       `bson:"oper_code"`
	CreateTime time.Time `bson:"create_time"`
}
