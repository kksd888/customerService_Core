package model

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var (
	Db *mgo.Database

	KfLastTimeChange = make(chan *Kf, 10)
)

func NewMongo() {
	session, err := mgo.Dial("172.16.14.52:27017")
	if err != nil {
		panic(err.Error())
	}
	Db = session.DB("test")

	go DbJob()
}

func DbJob() {
	kefuC := Db.C("kefu")
	for {
		k := <-KfLastTimeChange
		log.Printf("更新客服[%s]最后活动时间，%s", k.Id, k.UpdateTime)
		if err := kefuC.Update(bson.M{"id": k.Id}, bson.M{"$set": bson.M{"update_time": k.UpdateTime}}); err != nil {
			log.Printf("异步更新客服最后活动时间异常: %s", err.Error())
		}
	}
}
