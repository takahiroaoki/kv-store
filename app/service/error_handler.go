package service

import (
	"context"
	"fmt"

	"github.com/takahiroaoki/kv-store/app/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func handleError(ctx context.Context, appErr util.AppErr) error {
	if appErr == nil {
		return nil
	}

	switch appErr.LogLevel() {
	case util.LOG_LEVEL_INFO:
		util.InfoLogWithContext(ctx, appErr.Error())
	case util.LOG_LEVEL_WARN:
		util.WarnLogWithContext(ctx, appErr.Error())
	case util.LOG_LEVEL_ERROR:
		util.ErrorLogWithContext(ctx, appErr.Error())
	default:
		util.ErrorLogWithContext(ctx, fmt.Sprintf("error with no log level: %s", appErr.Error()))
	}

	var errMap = map[util.ErrorCause]error{
		util.CAUSE_UNDEFINED:        status.Error(codes.Unknown, "unknown error"),
		util.CAUSE_INVALID_ARGUMENT: status.Error(codes.InvalidArgument, "invalid argument"),
		util.CAUSE_NOT_FOUND:        status.Error(codes.NotFound, "data not found"),
		util.CAUSE_INTERNAL:         status.Error(codes.Internal, "internal error"),
	}

	switch appErr.Cause() {
	case util.CAUSE_UNDEFINED:
		util.ErrorLogWithContext(ctx, fmt.Sprintf("error with no cause: %s", appErr.Error()))
		return errMap[appErr.Cause()]
	default:
		return errMap[appErr.Cause()]
	}
}
