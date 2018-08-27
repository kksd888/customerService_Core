package logic

import (
	"github.com/satori/go.uuid"
	"time"
)

// 目前这个数据集写在内存中，之后考虑进行DB存储或者Redis存储
var RoomMap map[string]*Room

// 客户聊天会话数据集，Room对象维护实时的在线聊天信息、客服的在线状态。
// Room对象生命周期：
// * 由客户的访问进行初始化创建，并放入在线Room列表中
// * 客户和客服都可以更改各自对象中的在线状态
// * 由用户离线、客服终止来销毁Room对象
type Room struct {
	Id          uuid.UUID
	CustomerId  string
	KfId        string
	KfStatus    int
	CustomerMsg []string
	CreateTime  time.Time
}

func InitRoom(CustomerId string) *Room {
	uuid, _ := uuid.NewV4()
	return &Room{
		Id:         uuid,
		CustomerId: CustomerId,
		CreateTime: time.Time{},
	}
}

func (r *Room) Register() {
	_, ok := RoomMap[r.CustomerId]
	if !ok {
		RoomMap[r.CustomerId] = r
	}
}

func (r *Room) UnRegister(id uuid.UUID) {
	if r.Get(id) != nil {
		delete(RoomMap, r.CustomerId)
	}
}

func (r *Room) Get(id uuid.UUID) *Room {
	_, ok := RoomMap[r.CustomerId]
	if ok {
		return r
	} else {
		return nil
	}
}

func (r *Room) ChangeServerStatus(status int) {
	r.KfStatus = status
}
