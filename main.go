package main

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/controller/admin"
	"git.jsjit.cn/customerService/customerService_Core/controller/open"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"git.jsjit.cn/customerService/customerService_Core/wechat"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

var (
	wxContext *wechat.Wechat
)

func init() {
	wxContext = wechat.NewWechat(&wechat.Config{})
}

// @title 在线客服API文档
// @version 0.0.1
// @description  在线客服API文档的文档，接管了微信公众号聊天
// @BasePath /
func main() {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		log.Fatal("未找到配置文件conf.yaml", err)
	}
	config := common.GinConfig{}
	if err := yaml.Unmarshal(yamlFile, &config); err != nil {
		log.Fatal("配置文件格式错误", err)
	}

	gin.SetMode(config.RunModel)
	model.NewMongo(config.Mongodb)

	router := gin.Default()

	// cors规则配置
	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Authentication"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
		MaxAge:           30 * time.Minute,
	}))

	// 注册外部服务
	aiModule := handle.NewAiSemantic(config.AiSemantic)

	// Admin 注册控制器
	defaultController := admin.NewHealth()
	offlineReplyController := admin.NewOfflineReply()
	kfController := admin.NewKfServer()
	dialogController := admin.NewDialog(wxContext)
	weiXinController := admin.NewWeiXin(wxContext, aiModule)
	// OpenAPI注册控制器
	openController := open.NewOpen()
	openDialogController := open.NewDialog(aiModule)

	// 健康检查
	router.GET("/health", defaultController.Health)
	// 静态文件
	router.Static("/static", "./www")
	// 静态多媒体文件
	router.Static("/upload", "./upload")
	// 微信通信地址
	router.Any("/listen", weiXinController.Listen)
	// 客服登录操作
	router.POST("/adminGroup/login", kfController.LoginIn)

	// 后台Admin API路由 (授权保护)
	adminGroup := router.Group("/admin", handle.AdminOauthMiddleWare())
	{
		// 初始化
		adminGroup.GET("/init", defaultController.Init)

		// 待接入列表
		waitQueue := adminGroup.Group("/wait_queue")
		{
			waitQueue.GET("", dialogController.Queue)
			waitQueue.POST("/access", dialogController.Access)
		}

		// 会话操作
		dialog := adminGroup.Group("/dialog")
		{
			dialog.GET("", dialogController.List)
			dialog.GET("/:customerId/:page/:limit", dialogController.History)
			dialog.PUT("/ack", dialogController.Ack)
			dialog.POST("", dialogController.SendMessage)
		}

		// 客服操作
		kf := adminGroup.Group("/kf")
		{
			kf.GET("", kfController.Get)
			kf.PUT("/status", kfController.ChangeStatus)
		}

		// 设置操作
		setting := adminGroup.Group("/setting")
		{
			// 离线自动回复设置
			offlineReply := setting.Group("/offline_reply")
			{
				offlineReply.GET("", offlineReplyController.List)
				offlineReply.POST("", offlineReplyController.Create)
				offlineReply.PUT("/:replyId", offlineReplyController.Update)
				offlineReply.DELETE("/:replyId", offlineReplyController.Delete)
			}
		}
	}

	// API路由
	v1 := router.Group("/v1/app")
	{
		v1.POST("/access", openController.Access)

		// 授权保护
		d := v1.Group("/dialog", handle.OpenApiOauthMiddleWare())
		{
			d.GET("", openDialogController.Get)
			d.POST("", openDialogController.Create)
			d.GET("/history", openDialogController.History)
		}
	}

	go handle.Listen()

	// GO GO GO!!!
	router.Run(fmt.Sprintf(":%s", config.Port))
}
