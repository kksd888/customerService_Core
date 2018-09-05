package model

import (
	"gopkg.in/mgo.v2"
)

var (
	Db *mgo.Database
)

func NewMongo() {
	session, err := mgo.Dial("172.16.14.52:27017")
	if err != nil {
		panic(err.Error())
	}
	Db = session.DB("test")
}
