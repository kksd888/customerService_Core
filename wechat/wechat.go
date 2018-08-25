package wechat

import (
	"git.jsjit.cn/customerService/customerService_Core/wechat/cache"
	"git.jsjit.cn/customerService/customerService_Core/wechat/context"
	"git.jsjit.cn/customerService/customerService_Core/wechat/kf"
	"git.jsjit.cn/customerService/customerService_Core/wechat/server"
	"git.jsjit.cn/customerService/customerService_Core/wechat/user"
	"net/http"
	"sync"
)

// 微信上下文模型
type Wechat struct {
	Context *context.Context
}

// 配置
type Config struct {
	AppID          string
	AppSecret      string
	Token          string
	EncodingAESKey string
	Cache          cache.Cache
}

// 模型初始化
func NewWechat(cfg *Config) *Wechat {
	context := new(context.Context)
	copyConfigToContext(cfg, context)
	return &Wechat{context}
}

func copyConfigToContext(cfg *Config, context *context.Context) {
	context.AppID = cfg.AppID
	context.AppSecret = cfg.AppSecret
	context.Token = cfg.Token
	context.EncodingAESKey = cfg.EncodingAESKey
	context.Cache = cfg.Cache
	context.SetAccessTokenLock(new(sync.RWMutex))
	context.SetJsAPITicketLock(new(sync.RWMutex))
}

// GetServer 消息管理
func (wc *Wechat) GetServer(req *http.Request, writer http.ResponseWriter) *server.Server {
	wc.Context.Request = req
	wc.Context.Writer = writer
	return server.NewServer(wc.Context)
}

//GetAccessToken 获取access_token
func (wc *Wechat) GetAccessToken() (string, error) {
	return wc.Context.GetAccessToken()
}

// GetUser 用户管理接口
func (wc *Wechat) GetUser() *user.User {
	return user.NewUser(wc.Context)
}

//GetKf 客服管理接口
func (wc *Wechat) GetKf() *kf.Kf {
	return kf.NewCustomerServer(wc.Context)
}
