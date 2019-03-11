package message

import "encoding/xml"

// MsgType 基本消息类型
type MsgType string

// EventType 事件类型
type EventType string

const (
	//默认的成功消息，什么都不回复
	MsgTypeSuccess = "success"
	//MsgTypeText 表示文本消息
	MsgTypeText MsgType = "text"
	//MsgTypeImage 表示图片消息
	MsgTypeImage = "image"
	//MsgTypeVoice 表示语音消息
	MsgTypeVoice = "voice"
	//MsgTypeVideo 表示视频消息
	MsgTypeVideo = "video"
	//MsgTypeShortVideo 表示短视频消息[限接收]
	MsgTypeShortVideo = "shortvideo"
	//MsgTypeLocation 表示坐标消息[限接收]
	MsgTypeLocation = "location"
	//MsgTypeLink 表示链接消息[限接收]
	MsgTypeLink = "link"
	//MsgTypeMusic 表示音乐消息[限回复]
	MsgTypeMusic = "music"
	//MsgTypeNews 表示图文消息[限回复]
	MsgTypeNews = "news"
	//MsgTypeTransfer 表示消息消息转发到客服
	MsgTypeTransfer = "transfer_customer_service"
	//MsgTypeEvent 表示事件推送消息
	MsgTypeEvent = "event"
)

//MixMessage 存放所有微信发送过来的消息和事件
type MixMessage struct {
	CommonToken

	//基本消息
	MsgID        int64   `xml:"MsgId"`
	Content      string  `xml:"Content"`
	Recognition  string  `xml:"Recognition"`
	PicURL       string  `xml:"PicUrl"`
	MediaID      string  `xml:"MediaId"`
	Format       string  `xml:"Format"`
	ThumbMediaID string  `xml:"ThumbMediaId"`
	LocationX    float64 `xml:"Location_X"`
	LocationY    float64 `xml:"Location_Y"`
	Scale        float64 `xml:"Scale"`
	Label        string  `xml:"Label"`
	Title        string  `xml:"Title"`
	Description  string  `xml:"Description"`
	URL          string  `xml:"Url"`

	//事件相关
	Event       EventType `xml:"Event"`
	EventKey    string    `xml:"EventKey"`
	Ticket      string    `xml:"Ticket"`
	Latitude    string    `xml:"Latitude"`
	Longitude   string    `xml:"Longitude"`
	Precision   string    `xml:"Precision"`
	MenuID      string    `xml:"MenuId"`
	Status      string    `xml:"Status"`
	SessionFrom string    `xml:"SessionFrom"`

	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"`
		ScanResult string `xml:"ScanResult"`
	} `xml:"ScanCodeInfo"`

	SendPicsInfo struct {
		Count   int32      `xml:"Count"`
		PicList []EventPic `xml:"PicList>item"`
	} `xml:"SendPicsInfo"`

	SendLocationInfo struct {
		LocationX float64 `xml:"Location_X"`
		LocationY float64 `xml:"Location_Y"`
		Scale     float64 `xml:"Scale"`
		Label     string  `xml:"Label"`
		Poiname   string  `xml:"Poiname"`
	}
}

//EventPic 发图事件推送
type EventPic struct {
	PicMd5Sum string `xml:"PicMd5Sum"`
}

//EncryptedXMLMsg 安全模式下的消息体
type EncryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	ToUserName   string   `xml:"ToUserName" json:"ToUserName"`
	EncryptedMsg string   `xml:"Encrypt"    json:"Encrypt"`
}

//ResponseEncryptedXMLMsg 需要返回的消息体
type ResponseEncryptedXMLMsg struct {
	XMLName      struct{} `xml:"xml" json:"-"`
	EncryptedMsg string   `xml:"Encrypt"      json:"Encrypt"`
	MsgSignature string   `xml:"MsgSignature" json:"MsgSignature"`
	Timestamp    int64    `xml:"TimeStamp"    json:"TimeStamp"`
	Nonce        string   `xml:"Nonce"        json:"Nonce"`
}

// CommonToken 消息中通用的结构
type CommonToken struct {
	XMLName      xml.Name `xml:"xml" json:"xml_name"`
	ToUserName   string   `xml:"ToUserName" json:"to_user_name"`
	FromUserName string   `xml:"FromUserName" json:"from_user_name"`
	CreateTime   int64    `xml:"CreateTime" json:"create_time"`
	MsgType      MsgType  `xml:"MsgType" json:"msg_type"`
}

//SetToUserName set ToUserName
func (msg *CommonToken) SetToUserName(toUserName string) {
	msg.ToUserName = toUserName
}

//SetFromUserName set FromUserName
func (msg *CommonToken) SetFromUserName(fromUserName string) {
	msg.FromUserName = fromUserName
}

//SetCreateTime set createTime
func (msg *CommonToken) SetCreateTime(createTime int64) {
	msg.CreateTime = createTime
}

//SetMsgType set MsgType
func (msg *CommonToken) SetMsgType(msgType MsgType) {
	msg.MsgType = msgType
}
