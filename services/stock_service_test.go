package services

import (
	"net/http"
	"reflect"
	"stock-analyzer-wails/models"
	"testing"
	"time"
)

func TestStockService_GetStockByCode(t *testing.T) {
	type fields struct {
		baseURL string
		client  *http.Client
	}
	type args struct {
		code string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.StockData
		wantErr bool
	}{
		{
			name: "贵州茅台",
			fields: fields{
				baseURL: "http://78.push2.eastmoney.com/api/qt/clist/get",
				client: &http.Client{
					Timeout: 10 * time.Second,
				},
			},
			args: args{
				code: "603920",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StockService{
				//baseURL: tt.fields.baseURL,
				client: tt.fields.client,
			}
			got, err := s.GetStockByCode(tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetStockByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetStockByCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}
