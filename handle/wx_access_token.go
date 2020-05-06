package handle

import (
	"github.com/li-keli/go-tool/wechat/context"
)

//强制获取token
func GetQyAccessToken() (resAccessToken context.ResAccessToken, err error) {
	// 需要自行实现重新拉取Token的逻辑
	return
}
