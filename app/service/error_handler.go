package service

import (
	"context"
	"fmt"

	"github.com/takahiroaoki/go-libs/errorlibs"
	"github.com/takahiroaoki/kv-store/app/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(ctx context.Context, libErr errorlibs.Err) error {
	if libErr == nil {
		return nil
	}

	switch libErr.LogLevel() {
	case errorlibs.LOG_LEVEL_INFO:
		util.InfoLogWithContext(ctx, libErr.ErrorWithStackTrace())
	case errorlibs.LOG_LEVEL_WARN:
		util.WarnLogWithContext(ctx, libErr.ErrorWithStackTrace())
	case errorlibs.LOG_LEVEL_ERROR:
		util.ErrorLogWithContext(ctx, libErr.ErrorWithStackTrace())
	default:
		util.ErrorLogWithContext(ctx, fmt.Sprintf("error with no log level: %s", libErr.ErrorWithStackTrace()))
	}

	var errMap = map[errorlibs.ErrorCause]error{
		errorlibs.CAUSE_UNDEFINED:        status.Error(codes.Unknown, "unknown error"),
		errorlibs.CAUSE_INVALID_ARGUMENT: status.Error(codes.InvalidArgument, "invalid argument"),
		errorlibs.CAUSE_NOT_FOUND:        status.Error(codes.NotFound, "data not found"),
		errorlibs.CAUSE_INTERNAL:         status.Error(codes.Internal, "internal error"),
	}

	switch libErr.Cause() {
	case errorlibs.CAUSE_UNDEFINED:
		util.ErrorLogWithContext(ctx, fmt.Sprintf("error with no cause: %s", libErr.ErrorWithStackTrace()))
		return errMap[libErr.Cause()]
	default:
		return errMap[libErr.Cause()]
	}
}
