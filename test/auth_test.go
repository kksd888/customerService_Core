package test

import (
	"customerService_Core/handle"
	"fmt"
	"testing"
)

func Test_OpenApi(t *testing.T) {
	s, _ := handle.OpenMake2Auth("o1NTgjgYfl17vi5tSb5rPbN5MnhE")
	fmt.Println(s)
}
