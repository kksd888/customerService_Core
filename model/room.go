package model

import "time"

type RoomDb struct {
	Id         int       `gorm:"idgorm"`
	CustomerId string    `gorm:"customer_id"`
	KfId       int       `gorm:"kf_id"`
	KfStatus   int       `gorm:"kf_status"`
	CreateTime time.Time `gorm:"create_time"`
	UpdateTime time.Time `gorm:"update_time"`
}

func (RoomDb) TableName() string {
	return "dialog_room"
}

// 是否存在房间
func (r *RoomDb) IsExistByCustomerId(customerId string) bool {
	var count int
	firstErr := db.Where("customer_id=?", customerId).Count(&count)
	if firstErr.Error != nil {
		panic(firstErr.Error)
	}
	if count > 0 {
		return true
	} else {
		return false
	}
}
