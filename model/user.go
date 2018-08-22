package model

import "time"

// 用户模型
type User struct {
	Id           int
	OpenId       string
	NickName     string
	VisitCount   int
	CustomerType int
	CreateTime   time.Time
	UpdateTime   time.Time
}

func (u *User) insertOrUpdate(user User) error {
	return nil
}
