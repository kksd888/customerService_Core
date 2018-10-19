package test

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/controller/admin"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

var session *mgo.Session

func init() {
	session, _ = mgo.Dial("172.16.14.52:27018") // 测试数据库
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

func Test_Mongo_Update01(t *testing.T) {
	defer session.Close()

	roomCollection := session.DB("customer_service_db").C("room")
	_, updateErr := roomCollection.UpdateWithArrayFilters(
		bson.M{"room_customer.customer_id": "asdfasdf"},
		bson.M{"$set": bson.M{"room_messages.$[element].ack": true}},
		[]bson.M{
			{"element.oper_code": 2002},
		},
		true,
	)
	if updateErr != nil {
		log.Warn(updateErr)
	}

	ms := []model.Room{}
	roomCollection.Find(bson.M{"room_customer.customer_id": "asdfasdf", "room_messages.oper_code": common.MessageFromCustomer}).All(&ms)
	for _, v := range ms {
		for _, msg := range v.RoomMessages {
			fmt.Println(msg)
		}
	}
}

func Test_Mongo_Select(t *testing.T) {
	defer session.Close()

	var rooms []model.Room
	roomCollection := session.DB("customer_service_db").C("room")

	query := []bson.M{
		{
			"$match": bson.M{"room_kf.kf_id": "06f17d3d66194b24a72a3400db3fb9e9"},
		},
		{
			"$project": bson.M{
				"room_customer": 1,
				"room_messages": bson.M{
					"$filter": bson.M{
						"input": "$room_messages",
						"as":    "room_message",
						"cond": bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{"$$room_message.id", "24adb41d551642f3849aa8c476c49650"}},
								{"$eq": []interface{}{"$$room_message.ack", false}},
							},
						},
					},
				},
			},
		},
	}

	err := roomCollection.Pipe(query).All(&rooms)
	for _, v := range rooms {
		for _, msg := range v.RoomMessages {
			fmt.Println(msg)
		}
	}

	if err != nil {
		t.Fatal(err)
	}
}

func Test_Mongo_Select01(t *testing.T) {
	defer session.Close()

	var rooms []model.Room
	roomCollection := session.DB("customer_service_db").C("room")

	query := []bson.M{
		{
			"$match": bson.M{
				"$and": []bson.M{
					{"room_kf.kf_id": "06f17d3d66194b24a72a3400db3fb9e9"},
					{"room_messages.oper_code": bson.M{"$eq": common.MessageFromCustomer}},
					{"room_messages.ack": bson.M{"$eq": false}},
				},
			},
		},
		{
			"$project": bson.M{
				"room_customer": 1,
				"room_messages": bson.M{
					"$filter": bson.M{
						"input": "$room_messages",
						"as":    "room_message",
						"cond": bson.M{
							"$and": []bson.M{
								{"$eq": []interface{}{"$$room_message.oper_code", common.MessageFromCustomer}},
								{"$eq": []interface{}{"$$room_message.ack", false}},
							},
						},
					},
				},
			},
		},
	}

	if e := roomCollection.Pipe(query).All(&rooms); e != nil {
		log.Error(e)
	}

	newRoom := rooms[:0]
	for k, room := range rooms {
		if len(room.RoomMessages) > 1 {
			rooms = append(rooms[:k], rooms[(k+1):]...)
			newRoom = append(newRoom, room)
		}
	}
	for _, room := range newRoom {
		fmt.Println(room)
	}
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

	//var bsons []bson.M
	//roomCollection.Pipe([]bson.M{
	//	{
	//		"$match": bson.M{"room_kf.kf_id": "f24f257b370f4a6a9b703a35ea06f5b7"},
	//	},
	//	{
	//		"$project": bson.M{
	//			"room_messages": bson.M{"$slice": []interface{}{"$room_messages", -1}},
	//		},
	//	},
	//	{
	//		"$sort": bson.M{"room_messages.create_time": -1},
	//	},
	//	{
	//		"$limit": 100,
	//	},
	//}).All(&bsons)

	var bsons admin.RoomHistory
	roomCollection.Pipe([]bson.M{
		{
			"$match": bson.M{"room_customer.customer_nick_name": "只源有你"},
		},
		{
			"$unwind": "$room_messages",
		},
		{
			"$sort": bson.M{"room_messages.create_time": -1},
		},
		{
			"$skip": 0,
		},
		{
			"$limit": 10,
		},
		{
			"$group": bson.M{
				"_id":           "$_id",
				"room_messages": bson.M{"$push": "$room_messages"},
			},
		},
	}).One(&bsons)

	for _, v := range bsons.RoomMessages {
		fmt.Printf("%v \n", v)
	}
}

func Test_Times(t *testing.T) {
	fmt.Println(common.ToMd5("123JKD"))
	s, _ := admin.Make2Auth("5d893a28f68a4945a89a3f2db5f496f0")
	log.Println(s)
}

func Test_InitKf(t *testing.T) {
	defer session.Close()
	collection := session.DB("test").C("kefu")
	collection.Insert(&model.Kf{
		Id:         common.GetNewUUID(),
		JobNum:     "111",
		NickName:   "小金同学2",
		PassWord:   common.ToMd5("111"),
		HeadImgUrl: "http://thirdwx.qlogo.cn/mmopen/Q3auHgzwzM68w5nLXXsKOhFPqpB8wAyTz5TjXIHZ1ZfaroNrmPCjAJenrlrypP0XHl7WNf1vSW3AARJhNUryvoXTFsppf4ty3NicoA07kRQM/132",
		Type:       1,
		IsOnline:   false,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	})
}
