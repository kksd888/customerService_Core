package main

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"testing"
	"time"
)

var session *mgo.Session

func init() {
	session, _ = mgo.Dial("172.16.14.52:27017")
	mgo.SetDebug(true)
}

type User struct {
	Name string
	Age  int
	Msgs []UserMessage
}

type UserMessage struct {
	Id         int
	Msg        string
	CreateTime time.Time
}

func Test_Mongo_Insert(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("users")
	collection.Insert(&User{Name: "Admin", Age: 20, Msgs: []UserMessage{
		{Id: 1, Msg: "一个例子", CreateTime: time.Now()},
		{Id: 2, Msg: "第二个例子", CreateTime: time.Now()},
		{Id: 3, Msg: "第三个例子", CreateTime: time.Now()},
	}})
}

func Test_Mongo_Update(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("users")
	if e := collection.Update(bson.M{"age": 20}, bson.M{"$set": bson.M{"msgs.$[].msg": "修改成功2"}}); e != nil {
		t.Fatal(e.Error())
	}
}

func Test_Mongo_Select(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("users")
	query := collection.Find(bson.M{"msgs.id": 2})
	//if n, err := query.Count(); err != nil {
	//	t.Log(err)
	//} else {
	//	t.Log(n)
	//}

	iter := query.Iter()
	defer iter.Close()
	user := User{}
	for iter.Next(&user) {
		fmt.Printf("%v", user)
	}
}

func Test_InitKf(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("kefu")
	collection.Insert(&model.Kf{
		Id:         common.GetNewUUID(),
		TokenId:    "123",
		NickName:   "小金同学",
		HeadImgUrl: "http://thirdwx.qlogo.cn/mmopen/Q3auHgzwzM68w5nLXXsKOhFPqpB8wAyTz5TjXIHZ1ZfaroNrmPCjAJenrlrypP0XHl7WNf1vSW3AARJhNUryvoXTFsppf4ty3NicoA07kRQM/132",
		Type:       1,
		IsOnline:   false,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
}

func Test_Sclient(t *testing.T) {
	defer session.Close()
	roomCollection := session.DB("test").C("room")

	query := bson.M{
		"room_customer.customer_id": "ocnn-1PIPTsqqnRcVgUeIKCp2lKs",
	}
	changes := bson.M{
		"$push": bson.M{"room_messages": bson.M{"$each": []model.Message{
			{
				Id:         common.GetNewUUID(),
				Type:       "text",
				Msg:        "数组增量控制测试",
				MediaUrl:   "",
				OperCode:   common.MessageFromCustomer,
				CreateTime: time.Now(),
			},
		},
			"$slice": -10}},
	}
	if err := roomCollection.Update(query, changes); err != nil {
		log.Printf("异常消息：%s", err.Error())
	}
}

func Test_Sort(t *testing.T) {
	defer session.Close()
	roomCollection := session.DB("test").C("room")

	var bsons []bson.M
	roomCollection.Pipe([]bson.M{
		{
			"$match": bson.M{"room_kf.kf_id": "40e9146c-653c-4d75-8a6e-993ca1a0b34f"},
		},
		{
			"$project": bson.M{
				"room_messages": bson.M{"$slice": []interface{}{"$room_messages", 0, 1}},
			},
		},
		{
			"$sort": bson.M{"room_messages.create_time": -1},
		},
	}).All(&bsons)

	for _, v := range bsons {
		fmt.Printf("%#s \n", v["room_messages"])
	}
}
