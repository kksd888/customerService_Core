package test

import (
	"customerService_Core/handle"
	"fmt"
	"log"
	"testing"
	"time"
)

func TestAiSemantic_Dialogue(t *testing.T) {
	ai := handle.NewAiSemantic("http://localhost:5000/semantic")
	//s := ai.Dialogue("上海有贵宾厅吗？18888125808", "asdfasdfasdf")
	s := ai.Dialogue("你好，有人在吗？", "asdfasdfasdf")
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
