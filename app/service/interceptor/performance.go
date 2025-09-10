package interceptor

import (
	"context"
	"fmt"
	"time"

	"github.com/takahiroaoki/go-libs/timelibs"
	"github.com/takahiroaoki/kv-store/app/util"
	"google.golang.org/grpc"
)

func PerformanceLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (res any, err error) {
		reqAt := timelibs.Now()
		util.InfoLogWithContext(ctx, fmt.Sprintf("Request: %v", info.FullMethod))

		defer func() {
			latency := time.Since(reqAt.Time())
			util.PerfLogWithContext(ctx, fmt.Sprintf("Response: %v, Latency: %v", info.FullMethod, latency))
		}()

		res, err = handler(ctx, req)
		return
	}
}
