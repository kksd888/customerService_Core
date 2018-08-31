package logic

import "time"

// 目前这个数据集写在内存中，之后考虑进行DB存储或者Redis存储
var RoomMap = make(map[string]*Room)

// 客户聊天会话数据集，Room对象维护实时的在线聊天信息、客服的在线状态。
// Room对象生命周期：
// * 由客户的访问进行初始化创建，并放入在线Room列表中
// * 客户和客服都可以更改各自对象中的在线状态
// * 由用户离线、客服终止来销毁Room对象
type RoomContext interface {
	InitRoom(customerId string) (*Room, bool)
	UnRegister(customerId string)
	Update(room *Room) error
	GetWaitQueue() (waitQueueRooms []*Room, err error)
}

type Room struct {
	RoomCustomer
	RoomKf
	CreateTime time.Time
}
type RoomCustomer struct {
	CustomerId           string
	CustomerNickName     string
	CustomerHeadImgUrl   string
	CustomerPreviousKfId string
}
type RoomKf struct {
	KfId         int
	KfName       string
	KfHeadImgUrl string
	KfStatus     int
}
type RoomMessage struct {
	Uuid        string
	MessageType string
	Content     string
	CreateTime  time.Time
}
