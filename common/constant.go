package common

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
	"log"
	"net/http"
	"strings"
	"time"
)

// 客户来源
type CustomerSourceType string

const (
	AES_KEY = "80b11dc2dba242fd99b6bff28760c849" //AES加密的KEY

	KF_REPLY      = "您好，现在时段暂无人工客服为您服务，如您有任何问题可致电24小时服务热线4008101688"
	WELCOME_REPLY = "您好，请发送您要咨询的问题"
	LINE_UP_REPLY = "正有%d人排队，请稍后..."

	_              = iota // 客户类型
	NormalCustomer        // 普通客户
	VipCustomer           // VIP客户

	MessageFromSys      = 2001 // 系统信息
	MessageFromCustomer = 2002 // 客户发送的消息
	MessageFromKf       = 2003 // 客服发送的消息

	FromAPP    CustomerSourceType = "app"
	FromWeixin CustomerSourceType = "weixin"
)

// 异常码列表
const (
	StatusOk     = 100   // 正常
	Unknown      = 10000 // 系统未知异常
	ParameterBad = 10001 // 参数转换失败-请求参数异常
	HeaderBad    = 10002 // API请求头格式错误
	SignBad      = 10003 // 签名认证失败
)

const (
	//MsgTypeText 表示文本消息
	MsgTypeText MsgType = "text"
	//MsgTypeImage 表示图片消息
	MsgTypeImage = "image"
	//MsgTypeVoice 表示语音消息
	MsgTypeVoice = "voice"
)

var (
	// 本地时区
	LocalLocation, _ = time.LoadLocation("Local")
)

type MsgType string

type BaseOutput struct {
	Code   int         `json:"code"`
	ErrMsg string      `json:"err_msg"`
	Result interface{} `json:"result"`
}

// 成功返回
func ReturnSuccess(ctx *gin.Context, result interface{}) {
	ctx.JSON(http.StatusOK, &BaseOutput{
		Code:   StatusOk,
		Result: result,
	})
}

// 异常返回
func ReturnErr(ctx *gin.Context, err interface{}) {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, BaseOutput{
			Code:   Unknown,
			ErrMsg: err.(error).Error(),
			Result: nil,
		})
		log.Panicln(err)
	}
}

// 异常返回，自定义异常码
func ReturnErrCode(ctx *gin.Context, errCode int, err interface{}) {
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, BaseOutput{
			Code:   errCode,
			ErrMsg: err.(error).Error(),
			Result: nil,
		})
		log.Panicln(err)
	}
}

// 生成UUID
func GetNewUUID() string {
	uuids, _ := uuid.NewV4()
	return strings.Replace(uuids.String(), "-", "", -1)
}
