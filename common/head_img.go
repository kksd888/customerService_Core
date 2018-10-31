package common

import (
	"encoding/json"
	"git.jsjit.cn/customerService/customerService_Core/wechat/util"
	"math/rand"
	"time"
)

// 随机获取用户头像
func RandomHeadImg() string {
	defer func() {
		recover()
	}()

	var imgs baiduImg
	bytes, _ := util.HTTPGet("https://image.baidu.com/search/acjson?tn=resultjson_com&ipn=rj&fp=result&word=%E5%A4%B4%E5%83%8F+%E9%A3%8E%E6%99%AF")
	json.Unmarshal(bytes, &imgs)
	rand.Seed(time.Now().Unix())

	return imgs.Data[rand.Intn(len(imgs.Data))].MiddleURL
}

type baiduImg struct {
	Data []struct {
		MiddleURL string `json:"middleURL"`
	} `json:"data"`
}
