package services

//
//func TestStockService_GetStockByCode(t *testing.T) {
//	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
//		t.Skip("跳过需要外部网络的数据服务测试：设置 RUN_INTEGRATION_TESTS=1 可启用")
//	}
//
//	type fields struct {
//		baseURL string
//		client  *http.Client
//	}
//	type args struct {
//		code string
//	}
//	tests := []struct {
//		name    string
//		fields  fields
//		args    args
//		want    *models.StockData
//		wantErr bool
//	}{
//		{
//			name: "贵州茅台",
//			fields: fields{
//				baseURL: "http://78.push2.eastmoney.com/api/qt/clist/get",
//				client: &http.Client{
//					Timeout: 10 * time.Second,
//				},
//			},
//			args: args{
//				code: "603920",
//			},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			s := NewStockService()
//			s.client = tt.fields.client
//			got, err := s.GetStockByCode(tt.args.code)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("GetStockByCode() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetStockByCode() got = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}
