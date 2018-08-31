package logic

import (
	"encoding/base64"
	"errors"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/gin-gonic/gin/json"
	"log"
	"time"
)

type RoomMySql struct {
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
		if value.KfId == 0 {
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
