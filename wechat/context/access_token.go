package context

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"git.jsjit.cn/customerService/customerService_Core/common"
	"git.jsjit.cn/customerService/customerService_Core/wechat/util"
	"strconv"
)

const (
	//AccessTokenURL 获取access_token的接口
	AccessTokenURL = "https://api.weixin.qq.com/cgi-bin/token"
	//jsjAccessTokenURL = "http://172.16.5.63:9999/api/accesstoken?timestamp=%s&sign=%s"
	jsjAccessTokenURL = "http://wechat-mall.jsj.com.cn/api/accesstoken?timestamp=%s&sign=%s"
)

//ResAccessToken struct
type ResAccessToken struct {
	util.CommonError

	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

//SetAccessTokenLock 设置读写锁（一个appID一个读写锁）
func (ctx *Context) SetAccessTokenLock(l *sync.RWMutex) {
	ctx.accessTokenLock = l
}

//GetAccessToken 获取access_token
func (ctx *Context) GetAccessToken() (accessToken string, err error) {
	//ctx.accessTokenLock.Lock()
	//defer ctx.accessTokenLock.Unlock()
	//
	//accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.AppID)
	//val := ctx.Cache.Get(accessTokenCacheKey)
	//if val != nil {
	//	accessToken = val.(string)
	//	return
	//}

	//从微信服务器获取
	var resAccessToken ResAccessToken
	//resAccessToken, err = ctx.GetAccessTokenFromServer()
	resAccessToken, err = ctx.GetQyAccessTokenFromJsj()
	if err != nil {
		return
	}

	accessToken = resAccessToken.AccessToken
	return
}

//GetAccessTokenFromServer 强制从微信服务器获取token
func (ctx *Context) GetAccessTokenFromServer() (resAccessToken ResAccessToken, err error) {
	url := fmt.Sprintf("%s?grant_type=client_credential&appid=%s&secret=%s", AccessTokenURL, ctx.AppID, ctx.AppSecret)
	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrMsg != "" {
		err = fmt.Errorf("get access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	accessTokenCacheKey := fmt.Sprintf("access_token_%s", ctx.AppID)
	expires := resAccessToken.ExpiresIn - 1500
	err = ctx.Cache.Set(accessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	return
}

//GetQyAccessTokenFromJsj 强制从金色世纪获取token
func (ctx *Context) GetQyAccessTokenFromJsj() (resAccessToken ResAccessToken, err error) {
	unixTime := time.Now().Unix()
	sign := fmt.Sprint(unixTime, "jsjwechat*$(@^^^^)")

	url := fmt.Sprintf(jsjAccessTokenURL, strconv.FormatInt(unixTime, 10), common.ToMd5(sign))

	var body []byte
	body, err = util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get qy_access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	//qyAccessTokenCacheKey := fmt.Sprintf("qy_access_token_%s", ctx.AppID)
	//expires := resAccessToken.ExpiresIn - 1500
	//err = ctx.Cache.Set(qyAccessTokenCacheKey, resAccessToken.AccessToken, time.Duration(expires)*time.Second)
	return
}
