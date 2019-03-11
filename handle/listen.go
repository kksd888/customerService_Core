package handle

import (
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/globalsign/mgo/bson"
	"github.com/li-keli/go-tool/util/db_util"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	sigs = make(chan os.Signal, 1)
)

// 监听 优雅停服、客服下线，用户消息过期
func Listen() {

	// 跟新最后活动时间
	go func() {
		session := db_util.MongoDbSession.Copy()
		defer session.Close()

		kefuC := session.DB(common.DB_NAME).C("kefu")
		for {
			k := <-model.KfLastTimeChange
			//log.Printf("更新客服[%s]最后活动时间，%s", k.Id, k.UpdateTime)
			if err := kefuC.Update(bson.M{"id": k.Id}, bson.M{"$set": bson.M{"update_time": k.UpdateTime}}); err != nil {
				log.Printf("异步更新客服最后活动时间异常: %s", err.Error())
			}
		}
	}()

	// 优雅停机
	go func() {
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		// todo 所有客服强制下线，再议
		//var kefuC = model.Db.C("kefu")
		//kefuC.Update(bson.M{"is_online": true}, bson.M{"$set": bson.M{"is_online": false}})
		//log.Println("优雅停机，所有客服已经强制下线")
		os.Exit(0)
	}()

	// 客服超时下线
	go func() {
		session := db_util.MongoDbSession.Copy()
		defer session.Close()

		var kefuC = session.DB(common.DB_NAME).C("kefu")

		for {
			time.Sleep(time.Second * 1)

			var (
				duration, _ = time.ParseDuration("-10s")
				targetTime  = time.Now().Add(duration)
				allOffKefu  = []model.Kf{}
			)

			if err := kefuC.Find(bson.M{"is_online": true, "update_time": bson.M{"$lte": targetTime}}).All(&allOffKefu); err != nil {
				log.Printf(err.Error())
			}

			// 超时10秒则下线
			if len(allOffKefu) > 0 {
				for _, singeKf := range allOffKefu {
					if upErr := kefuC.Update(bson.M{"id": singeKf.Id}, bson.M{"$set": bson.M{"is_online": false, "update_time": time.Now()}}); upErr != nil {
						log.Printf("客服[%s]超时，下线异常: %s", singeKf.Id, upErr.Error())
					} else {
						log.Printf("客服[%s]超时，已经下线", singeKf.Id)
					}
				}
			}
		}
	}()
}
