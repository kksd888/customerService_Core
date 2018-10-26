package test

import (
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"log"
	"testing"
)

func TestAiSemantic_Dialogue(t *testing.T) {
	ai := handle.NewAiSemantic("http://172.16.14.55:20600/semantic")
	s := ai.Dialogue("珠海有贵宾厅吗？", "asdfasdfasdf")
	log.Println(s)
}
