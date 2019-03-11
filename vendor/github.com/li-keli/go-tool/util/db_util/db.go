package db_util

import (
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

var MongoDbSession *mgo.Session

// mongodb conn init
func NewMongo(url string) {
	session, err := mgo.Dial(url)
	if err != nil {
		logrus.Fatal("mongodb connection error: ", err, url)
	}
	session.SetMode(mgo.Monotonic, true)
	MongoDbSession = session
}
