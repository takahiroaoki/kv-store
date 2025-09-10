package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/go-libs/errorlibs"
	"github.com/takahiroaoki/kv-store/app/storage"
	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func TestKvServiceServer_Delete(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		req *pb.DeleteRequest
	}
	tests := []struct {
		name       string
		args       args
		setupMock  func(s *storage.MockStorage)
		assertion  assert.ErrorAssertionFunc
		want       *pb.DeleteResponse
		wantErrMsg string
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &pb.DeleteRequest{
					Key: "key",
				},
			},
			setupMock: func(s *storage.MockStorage) {
				s.EXPECT().DeleteKey(gomock.Any(), "key").Return(nil)
			},
			assertion: assert.NoError,
			want:      &pb.DeleteResponse{},
		},
		{
			name: "failure",
			args: args{
				ctx: context.Background(),
				req: &pb.DeleteRequest{
					Key: "key",
				},
			},
			setupMock: func(s *storage.MockStorage) {
				s.EXPECT().DeleteKey(gomock.Any(), "key").Return(errorlibs.NewErr(errors.New("error"), errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR))
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
			got, err := s.Delete(tt.args.ctx, tt.args.req)
			tt.assertion(t, err)
			if err != nil {
				assert.EqualError(t, err, tt.wantErrMsg)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
