package service

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/kv-store/app/model"
	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func Test_newKeyValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		req  *pb.SetRequest
		want model.KeyValue
	}{
		{
			name: "Success",
			req: &pb.SetRequest{
				Key:   "key",
				Value: "value",
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
			assert.Equal(t, tt.want, newKeyValue(tt.req))
		})
	}
}
