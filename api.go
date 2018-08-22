package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

// 轮询监听
func Listen(context *gin.Context) {
	context.JSON(200, gin.H{"status": "ok"})
}

// 长连接监听
func LongConn(context *gin.Context) {
	w := context.Writer
	r := context.Request

	var conn *websocket.Conn
	var err error
	conn, err = websocket.Upgrader{
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
		HandshakeTimeout: 5 * time.Second,
		CheckOrigin: func(r *http.Request) bool {
			// 取消ws跨域校验
			return true
		},
	}.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("websocket Upgrade err : %#v", err)
		return
	}

	for {
		t, reply, err := conn.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println(t, reply)

		// TODO 业务操作
	}
}
