package handle

import (
	"log"
	"testing"
)

func TestAiSemantic_Dialogue(t *testing.T) {
	ai := NewAiSemantic("http://172.16.14.55:20600/semantic")
	s := ai.Dialogue("珠海有贵宾厅吗？")
	log.Println(s)
}
