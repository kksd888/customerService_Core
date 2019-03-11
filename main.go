package main

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/controller/admin"
	"git.jsjit.cn/customerService/customerService_Core/controller/open"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"git.jsjit.cn/customerService/customerService_Core/model"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/li-keli/go-tool/wechat"
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
	// 加载配置
	common.NewGinConfig()
	// 数据库连接
	model.NewMongo()

	gin.SetMode(common.AppConfig.GoMode)
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
	aiModule := handle.NewAiSemantic(common.AppConfig.AiSemantic)

	// Admin 注册控制器
	adminController := admin.NewHealth()
	kfController := admin.NewKfServer()
	dialogController := admin.NewDialog(wxContext)
	statisticsController := admin.NewStatistics()
	weiXinController := admin.NewWeiXin(wxContext, aiModule)
	// OpenAPI注册控制器
	openController := open.NewOpen()
	openDialogController := open.NewDialog(aiModule)

	// 健康检查
	router.GET("/health", openController.Health)
	// 静态文件
	router.Static("/static", "./www")
	// 静态多媒体文件
	router.Static("/upload", "./upload")
	// 微信通信地址
	router.Any("/listen", weiXinController.Listen)
	// 客服登录操作
	router.POST("/admin/login", kfController.LoginIn)

	// 后台Admin API路由 (授权保护)
	adminGroup := router.Group("/admin", handle.AdminOauthMiddleWare())
	{
		// 初始化
		adminGroup.GET("/init", adminController.Init)

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

		// 统计操作
		statistics := adminGroup.Group("/statistics")
		{
			statistics.GET("/:starTime/:endTime/:page/:limit", statisticsController.Statistics)
		}

		// 客服操作
		kf := adminGroup.Group("/kf")
		{
			kf.GET("", kfController.Get)
			kf.PUT("/status", kfController.ChangeStatus)
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
	router.Run(fmt.Sprintf(":%s", common.AppConfig.Port))
}

// 统计数据
// 时间区间内，客服回复的信息总量
// db.getCollection('message').find({'oper_code':2003, 'create_time':{ "$gte" : ISODate("2018-10-01T00:00:00Z"), "$lt" : ISODate("2018-10-31T00:00:00Z")}}).count()
// 时间区间内，接待的客户数量
// db.getCollection('message').distinct('customer_id',{'oper_code':2003, 'create_time':{ "$gte" : ISODate("2018-10-01T00:00:00Z"), "$lt" : ISODate("2018-10-31T00:00:00Z")}})
// 员工登录的线上API文档
// http://memberapi.jsjinfo.cn/hosts/Jmember.aspx?help&m1=Member.API.API.JUser.LoginEmployee&rsdll=Member.API.dll&rqdll=Member.API.dll
