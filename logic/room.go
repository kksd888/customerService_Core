package logic

import (
	"encoding/base64"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/gin-gonic/gin/json"
	"github.com/satori/go.uuid"
	"log"
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
	RoomCustomer
	RoomKf
	CreateTime time.Time
}
type RoomCustomer struct {
	CustomerId           string
	CustomerNickName     string
	CustomerHeadImgUrl   string
	CustomerPreviousKfId string
	CustomerMsgs         []*RoomMessage
}
type RoomKf struct {
	KfId         string
	KfName       string
	KfHeadImgUrl string
	KfStatus     int
}

type RoomMessage struct {
	Uuid       uuid.UUID
	Content    string
	CreateTime time.Time
}

func InitRoom(customerId string) (*Room, bool) {
	r, ok := RoomMap[customerId]
	if ok {
		return r, false
	} else {
		newRoom := &Room{
			RoomCustomer: RoomCustomer{CustomerId: customerId},
			CreateTime:   time.Now(),
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
	uuids, _ := uuid.NewV4()
	r.CustomerMsgs = append(r.CustomerMsgs, &RoomMessage{
		Uuid:       uuids,
		Content:    msg,
		CreateTime: time.Now(),
	})
}

func (r *Room) ChangeServerStatus(status int) {
	r.KfStatus = status
}

func (r RoomKf) Make2Auth() (string, error) {
	bytes, err := json.Marshal(r)
	if err != nil {
		log.Printf("Make2Auth JSON序列化err：%v", err)
	}

	encrypt := common.AesEncrypt{}
	byteInfo, err := encrypt.Encrypt(bytes)
	if err != nil {
		log.Printf("common.NewGoAES() err：%v", err)
	}

	return base64.StdEncoding.EncodeToString(byteInfo), err
}

func KfAccess(customerIds []string, kfId RoomKf) {
	for roomKey, room := range RoomMap {
		for _, cIds := range customerIds {
			if cIds == roomKey {
				room.RoomKf = kfId
			}
		}
	}
}

func UpdateRoom(r *Room) (err error) {
	_, ok := RoomMap[r.CustomerId]
	if ok {
		RoomMap[r.CustomerId] = r
	} else {
		err = errors.New("查询的会话房间不存在")
	}

	return
}

func getRoomFromMaps(customerId string) *Room {
	r, ok := RoomMap[customerId]
	if ok {
		return r
	} else {
		return nil
	}
}

func GetWaitQueue() (waitQueueRooms []*Room, err error) {
	for _, value := range RoomMap {
		if value.KfId == "" {
			waitQueueRooms = append(waitQueueRooms, value)
		}
	}
	return
}

func PrintRoomMap() {
	for key, value := range RoomMap {
		log.Printf("RoomMap Key: %s ; Value : %#v", key, value)
	}
}
