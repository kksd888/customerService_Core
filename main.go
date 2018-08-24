package main

import (
	"git.jsjit.cn/customerService/customerService_Core/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	healthController := controller.InitHealth()
	dialogController := controller.InitDialog()
	serverController := controller.InitServer()
	offlineReplyController := controller.InitOfflineReply()

	// 定义路由
	v1 := r.Group("/v1")
	{
		// 健康检查
		v1.GET("/health", healthController.Health)

		// 会话操作
		dialog := v1.Group("/dialog")
		{
			dialog.GET("init", dialogController.DialogInit)
			dialog.GET("list", dialogController.List)
			dialog.POST("create", dialogController.Create)
			dialog.GET("user/:id/history", dialogController.History)
			dialog.POST("user/:id/message", dialogController.SendMessage)
			dialog.DELETE("user/:id/message", dialogController.RecallMessage)
		}

		// 客服操作
		server := v1.Group("/server")
		{
			server.GET(":id", serverController.Get)
			server.POST(":id/status", serverController.ChangeStatus)
		}

		// 设置操作
		setting := v1.Group("/setting")
		{
			// 离线自动回复设置
			offlineReply := setting.Group("offline_reply")
			{
				offlineReply.GET("", offlineReplyController.Get)
				offlineReply.POST("", offlineReplyController.Create)
				offlineReply.PUT(":id", offlineReplyController.Update)
				offlineReply.DELETE(":id", offlineReplyController.Delete)
			}
		}
	}

	r.Run(":5000")
}
