package model

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

var db *gorm.DB

func init() {
	mySqlDb, err := gorm.Open("mysql", "root:fK2g0Zx6@tcp(172.16.14.52:3306)/customer_service_bak?parseTime=true")
	//defer db.Close()

	//mySqlDb.LogMode(true)
	if err != nil {
		log.Printf("MySql数据库连接异常，%v", err)
	}
	db = mySqlDb
}
