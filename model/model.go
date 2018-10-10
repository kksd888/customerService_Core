package model

import (
	"gopkg.in/mgo.v2"
)

var (
	Db *mgo.Database
)

// mongodb项目中当做类似于Redis一类的直读缓存介质，用来维护数据的最终一致性
func NewMongo(conn string) {
	session, err := mgo.Dial(conn)
	if err != nil {
		panic(err.Error())
	}
	Db = session.DB("customer_service_db")
}
