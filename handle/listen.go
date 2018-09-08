package handle

import (
	"git.jsjit.cn/customerService/customerService_Core/model"
	"gopkg.in/mgo.v2/bson"
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
		var kefuC = model.Db.C("kefu")

		for {
			time.Sleep(time.Second * 1)

			var (
				duration, _ = time.ParseDuration("-1m")
				targetTime  = time.Now().Add(duration)
				allOffKefu  = []model.Kf{}
			)

			if err := kefuC.Find(bson.M{"is_online": true, "update_time": bson.M{"$lte": targetTime}}).All(&allOffKefu); err != nil {
				log.Printf(err.Error())
			}

			// 超时1分钟则下线
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
