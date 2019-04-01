package handle

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/li-keli/go-tool/util"
	"github.com/li-keli/go-tool/util/http_util"
	"github.com/li-keli/go-tool/wechat/context"
)

//GetQyAccessTokenFromJsj 强制从金色世纪获取token （生产环境的配置）
func GetQyAccessTokenFromJsj() (resAccessToken context.ResAccessToken, err error) {
	unixTime := time.Now().Unix()
	sign := fmt.Sprint(unixTime, "jsjwechat*$(@^^^^)")

	url := fmt.Sprintf("http://wechat-mall.jsj.com.cn/api/accesstoken?timestamp=%s&sign=%s", strconv.FormatInt(unixTime, 10), util.ToMd5(sign))

	var body []byte
	body, err = http_util.HTTPGet(url)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &resAccessToken)
	if err != nil {
		return
	}
	if resAccessToken.ErrCode != 0 {
		err = fmt.Errorf("get qy_access_token error : errcode=%v , errormsg=%v", resAccessToken.ErrCode, resAccessToken.ErrMsg)
		return
	}

	return
}
