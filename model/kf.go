package model

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"github.com/globalsign/mgo/bson"
	"github.com/li-keli/go-tool/util/db_util"
	"log"
	"time"
)

// 客服活动时间
var KfLastTimeChange = make(chan *Kf, 10)

// 客服模型
type Kf struct {
	Id         string    `json:"id" bson:"id"`
	JobNum     string    `json:"job_num" bson:"job_num"`
	NickName   string    `json:"nick_name" bson:"nick_name"`
	PassWord   string    `json:"-" bson:"pass_word"`
	Type       int       `json:"type" bson:"type"`
	HeadImgUrl string    `json:"head_img_url" bson:"head_img_url"`
	IsOnline   bool      `json:"is_online" bson:"is_online"`
	CreateTime time.Time `json:"create_time" bson:"create_time"`
	UpdateTime time.Time `json:"update_time" bson:"update_time"`
}

// 指定在线客服是否存在
func (k Kf) OnlineExist() bool {
	session := db_util.MongoDbSession.Copy()
	defer session.Close()

	var kefuC = session.DB(common.AppConfig.DbName).C("kefu")

	if count, err := kefuC.Find(bson.M{"id": k.Id, "is_online": true}).Count(); err != nil {
		log.Printf("model.Kf.Exist() is err :%s", err.Error())
		return false
	} else {
		if count > 0 {
			return true
		} else {
			return false
		}
	}
}

// 获取所有在线客服
func (k Kf) QueryOnlines() ([]*Kf, error) {
	session := db_util.MongoDbSession.Copy()
	defer session.Close()

	var (
		err     error
		onlines []*Kf
		kefuC   = session.DB(common.AppConfig.DbName).C("kefu")
	)
	if err = kefuC.Find(bson.M{"is_online": true}).All(&onlines); err != nil {
		log.Printf("model.QueryOnlines is err: %s", err.Error())
	}

	return onlines, err
}

// 修改客服在线状态
func (k Kf) ChangeStatus() (err error) {
	session := db_util.MongoDbSession.Copy()
	defer session.Close()

	kefuC := session.DB(common.AppConfig.DbName).C("kefu")

	if err = kefuC.Update(bson.M{"id": k.Id}, bson.M{"$set": bson.M{"is_online": k.IsOnline, "update_time": time.Now()}}); err != nil {
		log.Printf("model.ChangeStatus is err: %s", err.Error())
	}
	return
}
