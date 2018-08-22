package main

import (
	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	// 数据输出：计划支持轮询和长连接两种，轮询用于低版本游览器，高版本游览器使用长连接来提高性能
	// 处理轮询
	r.GET("/listen", Listen)

	// TODO 处理长连接
	r.GET("/ws")

	r.Run(":5000")
}
