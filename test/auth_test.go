package test

import (
	"customerService_Core/controller/admin"
	"customerService_Core/handle"
	"fmt"
	"testing"
)

func Test_AdminApi(t *testing.T) {
	s, _ := admin.Make2Auth("5d893a28f68a4945a89a3f2db5f496f0")
	fmt.Println(s)
}

func Test_OpenApi(t *testing.T) {
	s, _ := handle.OpenMake2Auth("o1NTgjgYfl17vi5tSb5rPbN5MnhE")
	fmt.Println(s)
}
