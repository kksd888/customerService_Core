package common

import (
	"github.com/satori/go.uuid"
	"strings"
)

const (
	AES_KEY  = "80b11dc2dba242fd99b6bff28760c849"                //AES加密的KEY
	KF_REPLY = "您好，现在时段暂无人工客服为您服务，如您有任何问题可致电24小时服务热线4008101688。" // 自动回复

	KF_ONLINE  = 0  // 客服在线
	KF_OFFLINE = -1 // 客服离线

	_              = iota // 客户类型
	NormalCustomer        // 普通客户
	VipCustomer           // VIP客户

	MessageFromCustomer = 2002 // 客户发送的消息
	MessageFromKf       = 2003 // 客服发送的消息
)

func GetNewUUID() string {
	uuids, _ := uuid.NewV4()
	return strings.Replace(uuids.String(), "-", "", -1)
}
