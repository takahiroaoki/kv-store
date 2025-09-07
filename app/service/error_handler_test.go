package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/kv-store/app/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_handleError(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		appErr util.AppErr
	}

	tests := []struct {
		name     string
		args     args
		expected error
	}{
		{
			name: "not found",
			args: args{
				ctx:    context.Background(),
				appErr: util.NewAppErrFromMsg("err", util.CAUSE_NOT_FOUND, util.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.NotFound, "data not found"),
		},
		{
			name: "invalid argument",
			args: args{
				ctx:    context.Background(),
				appErr: util.NewAppErrFromMsg("err", util.CAUSE_INVALID_ARGUMENT, util.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.InvalidArgument, "invalid argument"),
		},
		{
			name: "internal",
			args: args{
				ctx:    context.Background(),
				appErr: util.NewAppErrFromMsg("err", util.CAUSE_INTERNAL, util.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.Internal, "internal error"),
		},
		{
			name: "undefined",
			args: args{
				ctx:    context.Background(),
				appErr: util.NewAppErrFromMsg("err", util.CAUSE_UNDEFINED, util.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.Unknown, "unknown error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, handleError(tt.args.ctx, tt.args.appErr))
		})
	}
}
