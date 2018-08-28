package model

import (
	"log"
	"time"
)

// 用户模型
type Customer struct {
	Id           int       `gorm:"id"`
	OpenId       string    `gorm:"open_id"`
	NickName     string    `gorm:"nick_name"`
	Sex          int32     `gorm:"sex"`
	HeadImgUrl   string    `gorm:"head_img_url"`
	Address      string    `gorm:"address"`
	CustomerType int       `gorm:"customer_type"`
	CreateTime   time.Time `gorm:"create_time"`
	UpdateTime   time.Time `gorm:"update_time"`
}

// 新增或更新用户基础数据
func (customer Customer) InsertOrUpdate() (err error) {
	exec := db.Exec("replace into dic_customer (open_id, nick_name, customer_type, sex, head_img_url, address) values (?, ?, ?, ?, ?, ?);",
		customer.OpenId, customer.NickName, customer.CustomerType, customer.Sex, customer.HeadImgUrl, customer.Address)
	if exec.Error != nil {
		err = exec.Error
		log.Fatalf("Customer.InsertOrUpdate() is err => %#v", exec.Error)
	}
	return
}
