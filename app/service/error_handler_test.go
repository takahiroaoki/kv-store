package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/takahiroaoki/go-libs/errorlibs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_handleError(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx    context.Context
		libErr errorlibs.Err
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
				libErr: errorlibs.NewErrFromMsg("err", errorlibs.CAUSE_NOT_FOUND, errorlibs.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.NotFound, "data not found"),
		},
		{
			name: "invalid argument",
			args: args{
				ctx:    context.Background(),
				libErr: errorlibs.NewErrFromMsg("err", errorlibs.CAUSE_INVALID_ARGUMENT, errorlibs.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.InvalidArgument, "invalid argument"),
		},
		{
			name: "internal",
			args: args{
				ctx:    context.Background(),
				libErr: errorlibs.NewErrFromMsg("err", errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.Internal, "internal error"),
		},
		{
			name: "undefined",
			args: args{
				ctx:    context.Background(),
				libErr: errorlibs.NewErrFromMsg("err", errorlibs.CAUSE_UNDEFINED, errorlibs.LOG_LEVEL_INFO),
			},
			expected: status.Error(codes.Unknown, "unknown error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, handleError(tt.args.ctx, tt.args.libErr))
		})
	}
}
