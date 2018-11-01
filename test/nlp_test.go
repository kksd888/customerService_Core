package test

import (
	"fmt"
	"git.jsjit.cn/customerService/customerService_Core/handle"
	"log"
	"testing"
	"time"
)

func TestAiSemantic_Dialogue(t *testing.T) {
	ai := handle.NewAiSemantic("http://172.16.14.55:20700/semantic")
	s := ai.Dialogue("上海有贵宾厅吗？18888125808", "asdfasdfasdf")
	log.Println(s)
}

func Test_Date(t *testing.T) {
	now := time.Now()
	local1, _ := time.LoadLocation("")      // UTC 时区
	local2, _ := time.LoadLocation("Local") // 本地时区

	fmt.Println(now)
	fmt.Println(now.In(local1))
	fmt.Println(now.In(local2))
}
