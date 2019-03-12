package test

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/controller/admin"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	log "github.com/sirupsen/logrus"
	"testing"
	"time"
)

var session *mgo.Session

func init() {
	session, _ = mgo.Dial("172.16.14.52:27017") // 测试数据库，此处永远不准改成线上数据库
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

func Test_Mongo_Del01(t *testing.T) {
	defer session.Close()

	roomCollection := session.DB("customer_service_db").C("room")
	updateErr := roomCollection.Update(bson.M{"room_customer.customer_id": "C272E348914F4BF2A3DEB3B5262800E5"},
		bson.M{"$pull": bson.M{"room_messages": bson.M{"oper_code": 2002}}})
	if updateErr != nil {
		log.Error(updateErr)
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

func Test_Statistics(t *testing.T) {
	defer session.Close()

	//customerId := "o1NTgjoJdk49OC7pogphJzVpqz4s"
	//starTime := "2018-01-01"
	//endTime := "2018-12-31"
	//query := bson.M{"room_customer.customer_id": customerId, "room_messages.create_time": bson.M{"$gte": starTime, "$lte": endTime}}
	//roomCollection := session.DB("customer_service_db").C("room_messages")
	//kefuMessageCount, _ := roomCollection.Find(query).Count()
	//log.Println(kefuMessageCount)
	//t.Log(kefuMessageCount)
	//
	//redis := admin.NewStatistics()
	//context := &gin.Context{
	//	Params: nil,
	//}
	//redis.MessageCountByKf(context)

	var (
		kfid        = ""
		starTimeStr = "2018-09-07 00:00:00"
		endTimeStr  = "2018-09-08 00:00:00"
	)
	if kfid == "" {
		kfid = "06f17d3d66194b24a72a3400db3fb9e9"
	}

	//timestamp := time.Now().Unix()
	//tm := time.Unix(timestamp, 0)

	starTime, err := time.Parse("2006-01-02 15:04:05", starTimeStr)
	if err != nil {

	}
	endTime, err := time.Parse("2006-01-02 15:04:05", endTimeStr)
	if err != nil {

	}
	var (
		queryMessage = []bson.M{
			{
				"$match": bson.M{"kf_id": bson.M{"$ne": ""}, "create_time": bson.M{"$gte": starTime, "$lt": endTime}},
			},
			{"$lookup": bson.M{
				"from":         "kefu",
				"localField":   "kf_id",
				"foreignField": "id",
				"as":           "kefu",
			}},
			{
				"$unwind": bson.M{
					"path":                       "$kefu",
					"preserveNullAndEmptyArrays": true,
				},
			},
			{
				"$sort": bson.M{"kf_id": 1},
			},

			{
				"$group": bson.M{
					"_id":          "$kf_id",
					"kfId":         bson.M{"$first": "$kf_id"},
					"fkName":       bson.M{"$first": "$kefu.nick_name"},
					"messageCount": bson.M{"$sum": 1},
				},
			},
			{
				"$skip": (1 - 1) * 5,
			},
			{
				"$limit": 5,
			},
		}
		queryCustomer = []bson.M{
			{
				"$match": bson.M{"create_time": bson.M{"$gte": starTime, "$lt": endTime}},
			},
			{
				"$sort": bson.M{"kf_id": 1},
			},
			{
				"$group": bson.M{
					"_id":           bson.M{"kf_id": "$kf_id", "customer_id": "$customer_id"},
					"kfId":          bson.M{"$first": "$kf_id"},
					"customerId":    bson.M{"$first": "$customer_id"},
					"customerCount": bson.M{"$sum": 1},
				},
			},
		}

		roomCollection = session.DB("customer_service_db").C("message")
	)
	var messageByKf []map[string]interface{}
	if err := roomCollection.Pipe(queryMessage).All(&messageByKf); err != nil {
		log.Warn(err)
	}

	var customerByKf []bson.M
	if err := roomCollection.Pipe(queryCustomer).All(&customerByKf); err != nil {
		log.Warn(err)
	}

	count := 0
	for i := 0; i < len(messageByKf); i++ {
		kfId := messageByKf[i]["kfId"].(string)
		count = 0
		for j := 0; j < len(customerByKf); j++ {
			if kfId == customerByKf[j]["kfId"].(string) {
				count++
			}
			messageByKf[i]["customerCount"] = count
		}
	}
}

func Test_Login(t *testing.T) {
	defer session.Close()
	redis := admin.NewKfServer()
	context := &gin.Context{
		Params: nil,
	}
	redis.LoginIn(context)

	collection := session.DB("customer_service_db").C("kefu_copy1")

	//collection.Insert(&model.Kf{
	//	Id:         "6666",
	//	JobNum:     "kangyong",
	//	NickName:   "康勇",
	//	IsOnline:   true,
	//	CreateTime: time.Now(),
	//	UpdateTime: time.Now(),
	//	Type:       1,
	//	GroupName:  "投诉组",
	//})

	var kf model.Kf
	if err := collection.Find(bson.M{
		"job_num": "kangyong",
	}).One(&kf); err != nil {
		println(1)
	} else {
		println(2)
	}

	if err := collection.Find(bson.M{
		"job_num": "kangyong1",
	}).One(&kf); err != nil {
		if e := collection.Update(bson.M{"job_num": "6094"}, bson.M{"$set": bson.M{"is_online": true, "group_name": "投诉组"}}); e != nil {
			t.Fatal(e.Error())
		}
	} else {
		//collection.Insert(&model.Kf{
		//	Id:         "6666",
		//	JobNum:     "kangyong",
		//	NickName:   "康勇",
		//	IsOnline:   true,
		//	CreateTime: time.Now(),
		//	UpdateTime: time.Now(),
		//	Type:       1,
		//	GroupName:  "技术组",
		//})
	}
	println(kf.GroupName)

}

func Test_ChangeKf(t *testing.T) {
	defer session.Close()

	redis := admin.NewRoom()
	context := &gin.Context{
		Params: nil,
	}
	redis.ChangeKf(context)

	//kfCollection := session.DB("customer_service_db").C("kefu")
	//
	//kfOnline := []model.Kf{}
	//
	//if err := kfCollection.Find(bson.M{
	//	"group_name": "投诉组",
	//	"is_online":  false,
	//	"id":         bson.M{"$ne": "f4340e42fbd6e2fd9e6164033daf3194"},
	//}).All(&kfOnline); err != nil {
	//
	//} else {
	//	if len(kfOnline) > 0 {
	//		seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	//		kfId := kfOnline[seed.Intn(len(kfOnline)-1)].Id
	//		mesCollection := session.DB("customer_service_db").C("message_copy1")
	//		if e := mesCollection.Update(bson.M{"id": "33adcd0fd8a54f40a3e832146ef1ec81"}, bson.M{"$set": bson.M{"kf_id": kfId}}); e != nil {
	//			println(e)
	//		}
	//	} else {
	//		println(len(kfOnline))
	//	}
	//
	//}

}
