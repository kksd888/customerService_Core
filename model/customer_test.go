package model

import (
	"testing"
)

func TestCustomer_InsertOrUpdate(t *testing.T) {
	type fields struct {
		OpenId       string
		NickName     string
		VisitCount   int
		CustomerType int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			"新增测试",
			fields{
				OpenId:       "1234",
				NickName:     "测试用户",
				VisitCount:   0,
				CustomerType: 0,
			},
			false,
		},
		{
			"更新测试",
			fields{
				OpenId:       "1234",
				NickName:     "我改了一个名字",
				VisitCount:   1,
				CustomerType: 1,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			customer := Customer{
				OpenId:       tt.fields.OpenId,
				NickName:     tt.fields.NickName,
				CustomerType: tt.fields.CustomerType,
			}
			if err := customer.InsertOrUpdate(); (err != nil) != tt.wantErr {
				t.Errorf("Customer.Insert() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
