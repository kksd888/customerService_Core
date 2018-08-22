package main

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/net/websocket"
)

// 轮询监听
func Listen(context *gin.Context) {
	context.JSON(200, gin.H{"status": "ok"})
}

// 长连接监听
func LongConn(context *gin.Context) {
	handler := websocket.Handler()
	handler.ServeHTTP(context.Writer, context.Request)
}
