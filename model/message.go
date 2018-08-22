package model

import "time"

// 消息模型
type Message struct {
	Id         int
	UserId     int
	ServerId   int
	Msg        string
	MsgType    int
	OperCode   int
	CreateTime time.Time
	UpdateTime time.Time
}

func (m Message) InsertOrUpdate() error {
	return nil
}
