package model

import (
	"gopkg.in/mgo.v2"
)

type MongoDb struct {
	*mgo.Database
}

func NewMongo() *MongoDb {
	session, err := mgo.Dial("172.16.14.52:27017")
	if err != nil {
		panic(err.Error())
	}
	db := session.DB("test")
	return &MongoDb{db}
}
