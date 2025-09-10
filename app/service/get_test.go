package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/go-libs/errorlibs"
	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/storage"
	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func TestKvServiceServer_Get(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *pb.GetRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(s *storage.MockStorage)
		assertion  assert.ErrorAssertionFunc
		want       *pb.GetResponse
		wantErrMsg string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &pb.GetRequest{
					Key: "key",
				},
			},
			setupMock: func(s *storage.MockStorage) {
				s.EXPECT().GetByKey(gomock.Any(), "key").Return(model.KeyValue{Key: "key", Value: "value"}, nil)
			},
			assertion: assert.NoError,
			want: &pb.GetResponse{
				Value: "value",
			},
		},
		{
			name: "failure",
			args: args{
				ctx: context.Background(),
				req: &pb.GetRequest{
					Key: "key",
				},
			},
			setupMock: func(s *storage.MockStorage) {
				s.EXPECT().GetByKey(gomock.Any(), "key").Return(model.KeyValue{}, errorlibs.NewErr(errors.New("error"), errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR))
			},
			assertion:  assert.Error,
			want:       nil,
			wantErrMsg: "rpc error: code = Internal desc = internal error",
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
			got, err := s.Get(tt.args.ctx, tt.args.req)
			tt.assertion(t, err)
			if err != nil {
				assert.EqualError(t, err, tt.wantErrMsg)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
