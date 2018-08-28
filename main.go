package main

import (
	"git.jsjit.cn/customerService/customerService_Core/controller"
	_ "git.jsjit.cn/customerService/customerService_Core/docs"
	"git.jsjit.cn/customerService/customerService_Core/logic"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/cache"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

var (
	wxContext *wechat.Wechat
)

func init() {
	redis := cache.NewRedis(&cache.RedisOpts{
		Host: "localhost:6379",
	})

	//配置微信参数
	config := &wechat.Config{
		AppID:          "wx6cfceff5167a6007",
		AppSecret:      "1c1a365155e23b491f4878afbb87b918",
		Token:          "1603411701",
		EncodingAESKey: "fTrvMnac80fBHFP63KTLFZAhfxdSq7c126yftPw3HO1",
		Cache:          redis,
	}
	wxContext = wechat.NewWechat(config)
}

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	defaultController := controller.InitHealth()
	offlineReplyController := controller.InitOfflineReply()
	serverController := controller.InitKfServer()
	weiXinController := controller.InitWeiXin(wxContext, logic.RoomMap)
	dialogController := controller.InitDialog(wxContext, logic.RoomMap)

	// API路由
	v1 := router.Group("/v1")
	{
		// 健康检查
		v1.Any("/health", defaultController.Health)
		v1.GET("/init", defaultController.Init)

		// 会话操作
		dialog := v1.Group("/dialog")
		{
			dialog.GET(":dialogId/list", dialogController.List)
			dialog.POST("create", dialogController.Create)
		}

		// 访客数据操作
		customer := v1.Group("customer")
		{
			customer.GET(":id/history", dialogController.History)
			customer.POST(":id/message", dialogController.SendMessage)
			customer.DELETE(":id/message", dialogController.RecallMessage)
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

	// GO GO GO!!!
	router.Run(":5000")
}
