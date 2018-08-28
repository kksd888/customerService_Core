package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var MySqlDb *sql.DB

func init() {
	db, err := sql.Open("mysql", "root:fK2g0Zx6@tcp(172.16.14.52:3306)/customer_service?parseTime=true")
	//defer db.Close()

	if err != nil {
		log.Fatalf("MySql数据库连接异常，%v", err)
	}

	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)

	if err := db.Ping(); err != nil {
		log.Fatalf("MySql数据库通信异常，%v", err)
	}

	MySqlDb = db
}
