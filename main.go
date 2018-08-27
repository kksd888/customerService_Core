package main

import (
	"git.jsjit.cn/customerService/customerService_Core/controller"
	_ "git.jsjit.cn/customerService/customerService_Core/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
)

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	weiXinController := controller.InitWeiXin()
	defaultController := controller.InitHealth()
	dialogController := controller.InitDialog()
	serverController := controller.InitKfServer()
	offlineReplyController := controller.InitOfflineReply()

	// API路由
	v1 := router.Group("/v1")
	{
		// 健康检查
		v1.Any("/health", defaultController.Health)
		v1.GET("/init", defaultController.Init)

		// 会话操作
		dialog := v1.Group("/dialog")
		{
			dialog.GET("list", dialogController.List)
			dialog.POST("create", dialogController.Create)
			dialog.GET("customer/:id/history", dialogController.History)
			dialog.POST("customer/:id/message", dialogController.SendMessage)
			dialog.DELETE("customer/:id/message", dialogController.RecallMessage)
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
				offlineReply.GET("", offlineReplyController.List)
				offlineReply.POST("", offlineReplyController.Create)
				offlineReply.PUT(":id", offlineReplyController.Update)
				offlineReply.DELETE(":id", offlineReplyController.Delete)
			}
		}
	}

	// API文档地址
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 微信通信地址
	router.Any("/listen", weiXinController.Listen)

	go func() {
		for {
			wxMsg := <-controller.WxMsgQueue
			kf := controller.Wc.GetKf()
			if msgResponse, err := kf.SendTextMsg(wxMsg.FromUserName, "主动发送消息测试"); err != nil {
				log.Printf("%#v", msgResponse)
			}
		}
	}()

	// GO GO GO!!!
	router.Run(":5000")
}
