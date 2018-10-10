package common

import (
	"crypto/md5"
	"fmt"
)

func ToMd5(s string) string {
	data := []byte(s)
	has := md5.Sum(data)
	md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	return md5str1
}
