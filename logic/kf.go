package logic

import "git.jsjit.cn/customerService/customerService_Core/model"

var OnLineKfs = make(map[int]*model.Kf)

//type OnLineKf interface {
//	Get() (kf model.Kf, isOk bool)
//	Add(kf model.Kf)
//}

type OnLineKfMySql struct {
}

func (o *OnLineKfMySql) KfOnline(kf model.Kf) {
	OnLineKfs[kf.Id] = &kf
}

// 随机取出一个在线的客服
func GetOnlineKf() (kf model.Kf, isOk bool) {
	if len(OnLineKfs) == 0 {
		return model.Kf{}, false
	}

	for _, v := range OnLineKfs {
		return *v, true
	}

	return
}

// 客服上线
func AddOnlineKf(kf model.Kf) {
	OnLineKfs[kf.Id] = &kf
}
