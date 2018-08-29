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
	HeadImgUrl string    `gorm:"head_img_url"`
	CreateTime time.Time `gorm:"create_time"`
	UpdateTime time.Time `gorm:"update_time"`
}

func (Kf) TableName() string {
	return "dic_kf"
}

func (kf Kf) InsertOrUpdate() (err error) {
	return
}

func (kf *Kf) GetByTokenId(tokenId string) error {
	find := db.Where("token_id = ?", tokenId).First(&kf)
	return find.Error
}

func (kf *Kf) Get() error {
	find := db.Where("id = ? ", kf.Id).First(kf)
	return find.Error
}
