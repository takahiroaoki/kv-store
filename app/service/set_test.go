package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/storage"
	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func TestKvServiceServer_Set(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *pb.SetRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(s *storage.MockStorage)
		assertion  assert.ErrorAssertionFunc
		want       *pb.SetResponse
		wantErrMsg string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &pb.SetRequest{
					Key:   "key",
					Value: "value",
				},
			},
			setupMock: func(s *storage.MockStorage) {
				s.EXPECT().InsertKeyValue(gomock.Any(), model.KeyValue{
					Key:   "key",
					Value: "value",
				}).Return(nil)
			},
			assertion: assert.NoError,
			want:      &pb.SetResponse{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockStorage := storage.NewMockStorage(ctrl)
			tt.setupMock(mockStorage)
			s := &kvServiceServer{
				storage: mockStorage,
			}
			got, err := s.Set(tt.args.ctx, tt.args.req)
			tt.assertion(t, err)
			if err != nil {
				assert.EqualError(t, err, tt.wantErrMsg)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
