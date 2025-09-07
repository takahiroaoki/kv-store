package interceptor

import (
	"context"

	"github.com/google/uuid"
	"github.com/takahiroaoki/kv-store/app/util"
	"google.golang.org/grpc"
)

func SetContext() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		ctx = context.WithValue(ctx, util.REQUEST_ID, uuid.New())
		res, err = handler(ctx, req)
		return
	}
}
