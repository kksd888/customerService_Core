package main

import (
	"git.jsjit.cn/customerService/customerService_Core/controller"
	_ "git.jsjit.cn/customerService/customerService_Core/docs"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"git.jsjit.cn/customerService/customerService_Core/wechat/cache"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"time"
)

var (
	wxContext *wechat.Wechat
)

func init() {
	redis := cache.NewRedis(&cache.RedisOpts{
		Host: "172.16.7.20:6379",
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

// @title 在线客服API文档
// @version 0.0.1
// @description  在线客服API文档的文档，接管了微信公众号聊天
// @BasePath /
func main() {

	//gin.SetMode(gin.ReleaseMode)

	router := gin.Default()

	// CORS同源规则配置
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Authentication"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           12 * time.Hour,
	}))

	mongoDb := model.NewMongo()
	defaultController := controller.InitHealth(mongoDb)
	offlineReplyController := controller.InitOfflineReply()
	kfController := controller.InitKfServer(mongoDb)
	weiXinController := controller.InitWeiXin(wxContext, mongoDb)
	dialogController := controller.InitDialog(wxContext, mongoDb)
	//customerController := controller.InitCustomer(wxContext, logic.RoomMap)

	// API路由 (授权保护)
	v1 := router.Group("/v1", handle.OauthMiddleWare())
	{
		// 初始化
		v1.GET("/init", defaultController.Init)

		// 待接入列表
		waitQueue := v1.Group("/wait_queue")
		{
			waitQueue.GET("", dialogController.Queue)
			waitQueue.POST("/access", dialogController.Access)
		}

		// 会话操作
		dialog := v1.Group("/dialog")
		{
			dialog.GET("/", dialogController.List)
			dialog.POST("/", dialogController.SendMessage)
			dialog.PUT("/ack", dialogController.Ack)
			dialog.GET("/:customerId/:page/:limit", dialogController.History)
		}

		// 客服操作
		kf := v1.Group("/kf")
		{
			kf.GET("/", kfController.Get)
			kf.PUT("/status", kfController.ChangeStatus)
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
	router.POST("/login/:tokenId", kfController.LoginIn)
	//login.DELETE("/:tokenId", kfController.LoginOut)
	// 健康检查
	router.GET("/health", defaultController.Health)
	// API文档地址
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 微信通信地址
	router.Any("/listen", weiXinController.Listen)

	// GO GO GO!!!
	router.Run(":5000")
}
