package admin

import (
	"customerService_Core/common"
	"customerService_Core/handle"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/json"
	"github.com/gorilla/websocket"
	"github.com/li-keli/go-tool/util/mongo_util"
	"github.com/li-keli/mgo/bson"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	onLineKfs  = make(map[string]*websocket.Conn, 5)
	wsupgrader = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func WsHandler(ctx *gin.Context) {
	var (
		token = ctx.DefaultQuery("token", "")
		kfId  string
		err   error
	)

	if kfId, err = handle.AdminAuthToken2Model(token); err != nil {
		ctx.JSON(http.StatusUnauthorized, nil)
		return
	}

	wsConn(ctx.Writer, ctx.Request, kfId)
}

func wsConn(w http.ResponseWriter, r *http.Request, kfId string) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	var (
		conn *websocket.Conn
		err  error
	)

	// 握手连接
	if conn, err = wsupgrader.Upgrade(w, r, nil); err != nil {
		logrus.Error(err)
		return
	}

	go func(c *websocket.Conn, id string) {
		onLineKfs[id] = c
		for {
			messageType, p, err := c.ReadMessage()
			if err != nil {
				logrus.WithField("info", "websocket conn 异常").Error(err)
				kfLoginOut(id)
				return
			}

			if string(p) == "+" {
				b, _ := json.Marshal(WebSocketConnModel{Type: 0, Body: "-"})
				if err := c.WriteMessage(messageType, b); err != nil {
					logrus.Error(err)
					return
				}
			}
		}
	}(conn, kfId)

	// 确认登录状态
	_ = session.DB(common.AppConfig.DbName).C("kefu").
		Update(bson.M{"id": kfId}, bson.M{"$set": bson.M{"is_online": true, "update_time": time.Now()}})
}

// 客服下线
func kfLoginOut(kfId string) {
	session := mongo_util.GetMongoSession()
	defer session.Close()

	delete(onLineKfs, kfId)
	_ = session.DB(common.AppConfig.DbName).C("kefu").
		Update(bson.M{"id": kfId}, bson.M{"$set": bson.M{"is_online": false, "update_time": time.Now()}})
	logrus.Infoln("客服下线 => " + kfId)
}

// 通过websocket给在线客服发送消息
func SendMsgToOnlineKf(kfId string, body WebSocketConnModel) {

	logrus.Info("WebSocket发送消息 ", kfId, body)

	b, _ := json.Marshal(body)
	if conn, exist := onLineKfs[kfId]; exist {
		if err := conn.WriteMessage(1, b); err != nil {
			logrus.Error(err)
			return
		}
	} else {
		// 不存在的连接，则广播发待接入
		SendMsgRadio(body)
	}
}

// 通过WebSocket给所有客服广播
func SendMsgRadio(body WebSocketConnModel) {
	logrus.Info("WebSocket广播 ", body)

	b, _ := json.Marshal(body)
	for _, c := range onLineKfs {
		if err := c.WriteMessage(1, b); err != nil {
			logrus.Error(err)
		}
	}
}

// WebSocket通信模型
type WebSocketConnModel struct {
	Type int    `json:"Type"`
	Body string `json:"Body"`
}
