package storage

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/util"
)

func Test_newLogRow(t *testing.T) {
	type args struct {
		kv       model.KeyValue
		isDelete bool
	}
	tests := []struct {
		name      string
		args      args
		setupMock func(t *util.MockAppTime)
		want      logRow
	}{
		{
			name: "isDelete: false",
			setupMock: func(t *util.MockAppTime) {
				t.EXPECT().Format().Return("2024-01-01T00:00:00.000")
			},
			args: args{
				kv: model.KeyValue{
					Key:   "key",
					Value: "value",
				},
				isDelete: false,
			},
			want: logRow{
				key:       "key",
				value:     "value",
				delFlag:   "0",
				updatedAt: "2024-01-01T00:00:00.000",
			},
		},
		{
			name: "isDelete: true",
			setupMock: func(t *util.MockAppTime) {
				t.EXPECT().Format().Return("2024-01-01T00:00:00.000")
			},
			args: args{
				kv: model.KeyValue{
					Key:   "key",
					Value: "value",
				},
				isDelete: true,
			},
			want: logRow{
				key:       "key",
				value:     "value",
				delFlag:   "1",
				updatedAt: "2024-01-01T00:00:00.000",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockTime := util.NewMockAppTime(ctrl)

			originalFunc := util.Now
			defer func() { util.Now = originalFunc }()
			util.Now = func() util.AppTime {
				return mockTime
			}

			tt.setupMock(mockTime)
			assert.Equal(t, tt.want, newLogRow(tt.args.kv, tt.args.isDelete))
		})
	}
}
