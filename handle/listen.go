package handle

import (
	"git.jsjit.cn/customerService/customerService_Core/controller"
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
		// todo 所有客服强制下线 ？
		os.Exit(0)
	}()

	// 客服超时下线
	go func() {
		for {
			time.Sleep(time.Second * 1)
			for k := range controller.OnlineKfs {
				onLineKf := controller.OnlineKfs[k]
				duration := time.Now().Sub(onLineKf.LastTime)
				// 超时一小时则下线
				if duration.Hours() > 1 {
					delete(controller.OnlineKfs, k)
					if err := model.Db.C("room").Update(bson.M{"room_kf.kf_id": onLineKf.KfId}, bson.M{"$set": bson.M{"kf_status": "0"}}); err != nil {
						log.Printf("客服[%s]超时，下线异常: %s", onLineKf.KfId, err.Error())
					} else {
						log.Printf("客服[%s]超时，已经下线", onLineKf.KfId)
					}
				}
			}
		}
	}()
}
