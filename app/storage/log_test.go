package storage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	proctime "github.com/takahiroaoki/go-libs/time"
	"github.com/takahiroaoki/kv-store/app/model"
)

func Test_newLogRow(t *testing.T) {
	jst := time.FixedZone("JST", 9*60*60)
	proctime.SetLocation(jst)
	originalFunc := proctime.Now
	defer func() { proctime.Now = originalFunc }()
	proctime.Now = func() proctime.Time {
		return proctime.NewTime(time.Date(2024, 1, 1, 0, 0, 0, 0, jst))
	}
	type args struct {
		kv       model.KeyValue
		isDelete bool
	}
	tests := []struct {
		name string
		args args
		want logRow
	}{
		{
			name: "isDelete: false",
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
				updatedAt: "2024-01-01T00:00:00",
			},
		},
		{
			name: "isDelete: true",
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
				updatedAt: "2024-01-01T00:00:00",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, newLogRow(tt.args.kv, tt.args.isDelete))
		})
	}
}

func Test_newKeyValueFromLogRow(t *testing.T) {
	t.Parallel()
	type args struct {
		logRow logRow
	}
	tests := []struct {
		name string
		args args
		want model.KeyValue
	}{
		{
			name: "success",
			args: args{
				logRow: logRow{
					key:       "key",
					value:     "value",
					delFlag:   "0",
					updatedAt: "2024-01-01T00:00:00.000",
				},
			},
			want: model.KeyValue{
				Key:   "key",
				Value: "value",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, newKeyValueFromLogRow(tt.args.logRow))
		})
	}
}
