package model

import (
	"time"
)

// 客服模型
type Kf struct {
	Id         int       `gorm:"id"`
	TokenId    string    `gorm:"token_id"`
	NickName   string    `gorm:"nick_name"`
	Type       int       `gorm:"type"`
	CreateTime time.Time `gorm:"create_time"`
	UpdateTime time.Time `gorm:"update_time"`
}

func (kf Kf) InsertOrUpdate() (err error) {
	return
}

func (kf *Kf) GetByTokenId(tokenId string) error {
	find := db.Table("dic_kf").Where("token_id = ?", tokenId).First(kf)
	return find.Error
}
