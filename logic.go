package main

import (
	"github.com/satori/go.uuid"
	"time"
)

// 目前这个数据集写在内存中，之后考虑进行DB存储或者Redis存储
var RoomMap map[uuid.UUID]*Room

// 客户聊天会话数据集，Room对象维护实时的在线聊天信息、客服的在线状态。
// Room对象生命周期：
// * 由客户的访问进行初始化创建，并放入在线Room列表中
// * 客户和客服都可以更改各自对象中的在线状态
// * 由用户离线、客服终止来销毁Room对象
type Room struct {
	Id           uuid.UUID
	CustomerId   int
	ServerId     int
	ServerStatus int
	CreateTime   time.Time
}

func (r *Room) Register() {
	_, ok := RoomMap[r.Id]
	if !ok {
		RoomMap[r.Id] = r
	}
}

func (r *Room) UnRegister(id uuid.UUID) {
	if r.Get(id) != nil {
		delete(RoomMap, r.Id)
	}
}

func (r *Room) Get(id uuid.UUID) *Room {
	_, ok := RoomMap[r.Id]
	if ok {
		return r
	} else {
		return nil
	}
}

func (r *Room) ChangeServerStatus(status int) {
	r.ServerStatus = status
}
