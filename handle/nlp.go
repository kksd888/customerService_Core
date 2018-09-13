package handle

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/wechat/util"
	"log"
	"net/url"
)

// AI语义处理
type AiSemantic struct {
	hostUrl string
}

func NewAiSemantic(aiHost string) *AiSemantic {
	return &AiSemantic{hostUrl: aiHost}
}

func (ai *AiSemantic) Dialogue(msg string) string {
	msg = url.QueryEscape(msg)
	bytes, err := util.HTTPGet(fmt.Sprintf("%s?msg=%s", ai.hostUrl, msg))
	if err != nil {
		log.Printf("AiSemantic.Dialogue is err :%s", err.Error())
		return ""
	} else {
		return string(bytes)
	}
}
