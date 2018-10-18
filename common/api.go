package common

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
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
