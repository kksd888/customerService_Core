package model

import "time"

// 客服模型
type Server struct {
	Id         int
	TokenId    string
	NickName   string
	ServeCount int
	Type       int
	CreateTime time.Time
	UpdateTime time.Time
}

func (s Server) InsertOrUpdate() error {
	return nil
}
