package wechat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"time"
)

const (
	// 请求中的消息类型
	MsgTypeText       = "text"
	MsgTypeImage      = "image"
	MsgTypeVoice      = "voice"
	MsgTypeVideo      = "video"
	MsgTypeShortVideo = "shortvideo"
	MsgTypeLocation   = "location"
	MsgTypeLink       = "link"
	MsgTypeEvent      = "event"
	// 事件类型
	EventSubscribe   = "subscribe"
	EventUnsubscribe = "unsubscribe"
	EventScan        = "SCAN"
	EventLocation    = "LOCATION"
	EventClick       = "CLICK"
	EventView        = "VIEW"
	// 多媒体类型
	MediaTypeImage = "image"
	MediaTypeVoice = "voice"
	MediaTypeVideo = "video"
	MediaTypeThumb = "thumb"
	// 环境
	UrlPrefix      = "https://api.weixin.qq.com/cgi-bin/"
	MediaUrlPrefix = "http://file.api.weixin.qq.com/cgi-bin/media/"
	retryNum       = 3
)

type WeChatService struct {
	Request     Request
	AccessToken AccessToken
}

func New(token, appId, appSecret string) *WeChatService {
	return &WeChatService{
		Request:     Request{Token: token},
		AccessToken: AccessToken{AppId: appId, AppSecret: appSecret},
	}
}

type msgHeader struct {
	XMLName      xml.Name `xml:"xml" json:"-"`
	ToUserName   string   `json:"touser"`
	FromUserName string   `json:"-"`
	CreateTime   int64    `json:"-"`
	MsgType      string   `json:"msgtype"`
}

type textMsg struct {
	msgHeader
	Content string `json:"-"`
	Text    struct {
		Content string `xml:"-" json:"content"`
	} `xml:"-" json:"text"`
}

type imageMsg struct {
	msgHeader
	Image struct {
		MediaId string `json:"media_id"`
	} `json:"image"`
}

type voiceMsg struct {
	msgHeader
	Voice struct {
		MediaId string `json:"media_id"`
	} `json:"voice"`
}

type videoMsg struct {
	msgHeader
	Video *Video `json:"video"`
}

type musicMsg struct {
	msgHeader
	Music *Music `json:"music"`
}

type newsMsg struct {
	msgHeader
	ArticleCount int `json:"-"`
	Articles     struct {
		Item *[]Article `xml:"item" json:"articles"`
	} `json:"news"`
}

type newsGroupMsg struct {
	// 群发图文消息

	Filter struct {
		IsToAll bool   `json:"is_to_all"`
		GroupId string `json:"group_id"`
	} `json:"filter"`

	Mpnews struct {
		MediaId string `json:"media_id"`
	} `json:"mpnews"`

	MsgType string `json:"msgtype"`
}

type Video struct {
	MediaId     string `json:"media_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Music struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	MusicUrl     string `json:"musicurl"`
	HQMusicUrl   string `json:"hqmusicurl"`
	ThumbMediaId string `json:"thumb_media_id"`
}

type Article struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PicUrl      string `json:"picurl"`
	Url         string `json:"url"`
}

// 回复文本
func (wcServer *WeChatService) ReplyTextMsg(rw http.ResponseWriter, content string) error {
	var msg textMsg
	msg.MsgType = "text"
	msg.Content = content
	return wcServer.replyMsg(rw, &msg)
}

// 回复图片
func (wcServer *WeChatService) ReplyImageMsg(rw http.ResponseWriter, mediaId string) error {
	var msg imageMsg
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	return wcServer.replyMsg(rw, &msg)
}

// 回复语音
func (wcServer *WeChatService) ReplyVoiceMsg(rw http.ResponseWriter, mediaId string) error {
	var msg voiceMsg
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	return wcServer.replyMsg(rw, &msg)
}

// 回复视频
func (wcServer *WeChatService) ReplyVideoMsg(rw http.ResponseWriter, video *Video) error {
	var msg videoMsg
	msg.MsgType = "video"
	msg.Video = video
	return wcServer.replyMsg(rw, &msg)
}

// 回复音乐
func (wcServer *WeChatService) ReplyMusicMsg(rw http.ResponseWriter, music *Music) error {
	var msg musicMsg
	msg.MsgType = "music"
	msg.Music = music
	return wcServer.replyMsg(rw, &msg)
}

// 回复图文消息
func (wcServer *WeChatService) ReplyNewsMsg(rw http.ResponseWriter, articles *[]Article) error {
	var msg newsMsg
	msg.MsgType = "news"
	msg.ArticleCount = len(*articles)
	msg.Articles.Item = articles
	return wcServer.replyMsg(rw, &msg)
}

func (wcServer *WeChatService) replyMsg(rw http.ResponseWriter, msg interface{}) error {
	v := reflect.ValueOf(msg).Elem()
	v.FieldByName("ToUserName").SetString(wcServer.Request.FromUserName)
	v.FieldByName("FromUserName").SetString(wcServer.Request.ToUserName)
	v.FieldByName("CreateTime").SetInt(time.Now().Unix())
	data, err := xml.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := rw.Write(data); err != nil {
		return err
	}
	return nil
}

// 发送文字
func (wcServer *WeChatService) SendTextMsg(touser string, content string) error {
	var msg textMsg
	msg.MsgType = "text"
	msg.Text.Content = content
	return wcServer.sendMsg(touser, &msg)
}

// 发送图片
func (wcServer *WeChatService) SendImageMsg(touser string, mediaId string) error {
	var msg imageMsg
	msg.MsgType = "image"
	msg.Image.MediaId = mediaId
	return wcServer.sendMsg(touser, &msg)
}

// 发送语音
func (wcServer *WeChatService) SendVoiceMsg(touser string, mediaId string) error {
	var msg voiceMsg
	msg.MsgType = "voice"
	msg.Voice.MediaId = mediaId
	return wcServer.sendMsg(touser, &msg)
}

// 发送视频
func (wcServer *WeChatService) SendVideoMsg(touser string, video *Video) error {
	var msg videoMsg
	msg.MsgType = "video"
	msg.Video = video
	return wcServer.sendMsg(touser, &msg)
}

// 发送音乐
func (wcServer *WeChatService) SendMusicMsg(touser string, music *Music) error {
	var msg musicMsg
	msg.MsgType = "music"
	msg.Music = music
	return wcServer.sendMsg(touser, &msg)
}

// 发送图文
func (wcServer *WeChatService) SendNewsMsg(touser string, articles *[]Article) error {
	var msg newsMsg
	msg.MsgType = "news"
	msg.Articles.Item = articles
	return wcServer.sendMsg(touser, &msg)
}

func (wcServer *WeChatService) sendMsg(touser string, msg interface{}) error {
	v := reflect.ValueOf(msg).Elem()
	v.FieldByName("ToUserName").SetString(touser)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%smessage/custom/send?access_token=", UrlPrefix)
	buf := bytes.NewBuffer(data)
	// retry
	for i := 0; i < retryNum; i++ {
		token, err := wcServer.AccessToken.Fresh()
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		if _, err := post(url+token, "text/plain", buf); err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// 向全部用户群发图文消息
func (wcServer *WeChatService) sendNewsToALl(mediaId string) error {
	var news newsGroupMsg
	news.MsgType = "mpnews"
	news.Filter.IsToAll = true
	news.Mpnews.MediaId = mediaId
	return wcServer.sendGroupMsg(news)
}

// 向特定GroupId用户群发图文消息
func (wcServer *WeChatService) sendNewsToGroup(groupId string, mediaId string) error {
	var news newsGroupMsg
	news.MsgType = "mpnews"
	news.Filter.IsToAll = false
	news.Filter.GroupId = groupId
	news.Mpnews.MediaId = mediaId
	return wcServer.sendGroupMsg(news)
}

// 群发消息
func (wcServer *WeChatService) sendGroupMsg(msg interface{}) error {
	url := fmt.Sprintf("%smessage/mass/sendall?access_token=", UrlPrefix)
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(data)
	// retry
	for i := 0; i < retryNum; i++ {
		token, err := wcServer.AccessToken.Fresh()
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		if _, err := post(url+token, "text/plain", buf); err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// 下载多媒体文件
func (wcServer *WeChatService) DownloadMediaFile(mediaId, fileName string) error {
	url := fmt.Sprintf("%sget?media_id=%s&access_token=", MediaUrlPrefix, mediaId)
	// retry
	for i := 0; i < retryNum; i++ {
		token, err := wcServer.AccessToken.Fresh()
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		resp, err := http.Get(url + token)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		// json
		if resp.Header.Get("Content-Type") == "text/plain" {
			var rtn response
			if err := json.Unmarshal(data, &rtn); err != nil {
				if i < retryNum-1 {
					continue
				}
				return err
			}
			if i < retryNum-1 {
				continue
			}
			return errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
		}
		// media
		f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		defer f.Close()

		if _, err := f.Write(data); err != nil {
			if i < retryNum-1 {
				continue
			}
			return err
		}
		break // success
	}
	return nil
}

// 上传多媒体文件
func (wcServer *WeChatService) UploadMediaFile(mediaType, fileName string) (string, error) {
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)
	defer bw.Close()
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	fw, err := bw.CreateFormFile("filename", f.Name())
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(fw, f); err != nil {
		return "", err
	}
	f.Close()
	bw.Close()
	url := fmt.Sprintf("%supload?type=%s&access_token=", MediaUrlPrefix, mediaType)
	mime := bw.FormDataContentType()
	mediaId := ""
	// retry
	for i := 0; i < retryNum; i++ {
		token, err := wcServer.AccessToken.Fresh()
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return "", err
		}
		rtn, err := post(url+token, mime, &buf)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return "", err
		}
		mediaId = rtn.MediaId
		break // success
	}
	return mediaId, nil
}

type UserInfo struct {
	Subscribe     int64  `json:"subscribe"`
	Openid        string `json:"openid"`
	Nickname      string `json:"nickname"`
	Sex           int64  `json:"sex"`
	Language      string `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	Headimgurl    string `json:"headimgurl"`
	SubscribeTime int64  `json:"subscribe_time"`
}

// 获取用户信息
func (wcServer *WeChatService) GetUserInfo(openId string) (UserInfo, error) {
	var uinf UserInfo
	url := fmt.Sprintf("%suser/info?lang=zh_CN&openid=%s&access_token=", UrlPrefix, openId)
	// retry
	for i := 0; i < retryNum; i++ {
		token, err := wcServer.AccessToken.Fresh()
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		resp, err := http.Get(url + token)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		// has error?
		var rtn response
		if err := json.Unmarshal(data, &rtn); err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		// yes
		if rtn.ErrCode != 0 {
			if i < retryNum-1 {
				continue
			}
			return uinf, errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
		}
		// no
		if err := json.Unmarshal(data, &uinf); err != nil {
			if i < retryNum-1 {
				continue
			}
			return uinf, err
		}
		break // success
	}
	return uinf, nil
}
