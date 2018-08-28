package model

import (
	"log"
	"time"
)

// 用户模型
type Customer struct {
	Id           int
	OpenId       string
	NickName     string
	Sex          int32
	HeadImgUrl   string
	Address      string
	CustomerType int
	CreateTime   time.Time
	UpdateTime   time.Time
}

// 新增或更新用户基础数据
func (customer Customer) InsertOrUpdate() (err error) {
	_, err = MySqlDb.Exec("replace into dic_customer (open_id, nick_name, customer_type, sex, head_img_url, address) values (?, ?, ?, ?, ?, ?);",
		customer.OpenId, customer.NickName, customer.CustomerType, customer.Sex, customer.HeadImgUrl, customer.Address)
	if err != nil {
		log.Fatalf("Customer.InsertOrUpdate() is err => %#v", err)
	}
	return
}
