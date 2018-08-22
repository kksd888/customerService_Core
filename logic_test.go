package main

import (
	"testing"
	"time"

	"github.com/satori/go.uuid"
)

func TestRoom_getRoom(t *testing.T) {
	type fields struct {
		Id           uuid.UUID
		CustomerId   int
		ServerId     int
		ServerStatus int
		CreateTime   time.Time
	}

	tests := []struct {
		name   string
		fields fields
		args   uuid.UUID
	}{
		{
			name:   "case1",
			fields: fields{},
			args: func() uuid.UUID {
				u, _ := uuid.NewV4()
				return u
			}(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Room{
				Id:           tt.fields.Id,
				CustomerId:   tt.fields.CustomerId,
				ServerId:     tt.fields.ServerId,
				ServerStatus: tt.fields.ServerStatus,
				CreateTime:   tt.fields.CreateTime,
			}
			room := r.Get(tt.args)
			if room == nil {
				t.Errorf("检索room列表，未找到")
			}
		})
	}
}
