package handle

import (
	"log"
	"testing"
)

func TestAiSemantic_Dialogue(t *testing.T) {
	ai := NewAiSemantic("http://127.0.0.1:5000/semantic")
	s := ai.Dialogue("珠海有贵宾厅吗？")
	log.Println(s)
}
