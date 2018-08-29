package main

import (
	"git.jsjit.cn/customerService/customerService_Core/controller"
	_ "git.jsjit.cn/customerService/customerService_Core/docs"
	"git.jsjit.cn/customerService/customerService_Core/handle"
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

	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(cors.Default())

	defaultController := controller.InitHealth()
	offlineReplyController := controller.InitOfflineReply()
	kfController := controller.InitKfServer()
	weiXinController := controller.InitWeiXin(wxContext, logic.RoomMap)
	dialogController := controller.InitDialog(wxContext, logic.RoomMap)
	customerController := controller.InitCustomer(wxContext, logic.RoomMap)

	// API路由 (授权保护)
	v1 := router.Group("/v1", handle.OauthMiddleWare())
	{
		// 初始化
		v1.GET("/init", defaultController.Init)

		// 待接入列表
		waitQueue := v1.Group("/wait_queue")
		{
			waitQueue.GET("", customerController.Queue)
		}

		// 会话操作
		dialog := v1.Group("/dialog")
		{
			dialog.GET("/:kfId/list", dialogController.List)
			dialog.POST("/access", dialogController.Access)
		}

		// 客户数据
		customer := v1.Group("/customer")
		{
			customer.GET("/:customerId/history/", customerController.History)
			customer.POST("/:customerId/message/", customerController.SendMessage)
		}

		// 客服操作
		kf := v1.Group("/kf")
		{
			kf.GET("/:kfId", kfController.Get)
			kf.POST("/:kfId/status", kfController.ChangeStatus)
		}

		// 设置操作
		setting := v1.Group("/setting")
		{
			// 离线自动回复设置
			offlineReply := setting.Group("/offline_reply")
			{
				offlineReply.GET("", offlineReplyController.List)
				offlineReply.POST("", offlineReplyController.Create)
				offlineReply.PUT("/:replyId/", offlineReplyController.Update)
				offlineReply.DELETE("/:replyId/", offlineReplyController.Delete)
			}
		}
	}

	// 客服登录操作
	login := router.Group("/v1/login")
	login.POST("/:tokenId", kfController.LoginIn)
	//login.DELETE("/:tokenId", kfController.LoginOut)
	// 健康检查
	router.Any("/health", defaultController.Health)
	// API文档地址
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 微信通信地址
	router.Any("/listen", weiXinController.Listen)

	// GO GO GO!!!
	router.Run(":5000")
}
