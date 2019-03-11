package model

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/globalsign/mgo"
)

var Db *mgo.Database

// mongodb项目中当做类似于Redis一类的直读缓存介质，用来维护数据的最终一致性
func NewMongo() {
	session, err := mgo.Dial(common.AppConfig.Mongodb)
	if err != nil {
		panic(err.Error())
	}
	Db = session.DB("customer_service_db")
}
