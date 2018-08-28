package logic

import (
	"time"
)

// 目前这个数据集写在内存中，之后考虑进行DB存储或者Redis存储
var RoomMap = make(map[string]*Room)

// 客户聊天会话数据集，Room对象维护实时的在线聊天信息、客服的在线状态。
// Room对象生命周期：
// * 由客户的访问进行初始化创建，并放入在线Room列表中
// * 客户和客服都可以更改各自对象中的在线状态
// * 由用户离线、客服终止来销毁Room对象
type Room struct {
	CustomerId  string
	KfId        string
	KfStatus    int
	CustomerMsg []string
	CreateTime  time.Time
}

func InitRoom(customerId string) (*Room, bool) {
	r, ok := RoomMap[customerId]
	if ok {
		return r, false
	} else {
		newRoom := &Room{
			CustomerId: customerId,
			CreateTime: time.Time{},
		}

		RoomMap[customerId] = newRoom
		return newRoom, true
	}
}

func (r *Room) UnRegister(customerId string) {
	if getRoomFromMaps(customerId) != nil {
		delete(RoomMap, r.CustomerId)
	}
}

func (r *Room) AddMessage(msg string) {
	r.CustomerMsg = append(r.CustomerMsg, msg)
}

func (r *Room) ChangeServerStatus(status int) {
	r.KfStatus = status
}

func getRoomFromMaps(customerId string) *Room {
	r, ok := RoomMap[customerId]
	if ok {
		return r
	} else {
		return nil
	}
}
