package wechat

// 微信消息处理器，处理和微信服务器的交互
type WeChat struct {
	token string
}

// 获取微信授权
func (wc *WeChat) getAccessToken() {
}

// 接受微信发送过来的数据
func (wc *WeChat) receiveWeChatMsg() {
}

// 发送消息给微信
func (wc *WeChat) sendWeChatMsg() {
}
