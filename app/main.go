package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/takahiroaoki/go-libs/timelibs"
	"github.com/takahiroaoki/kv-store/app/config"
	"github.com/takahiroaoki/kv-store/app/service"
	"github.com/takahiroaoki/kv-store/app/storage"
	"github.com/takahiroaoki/kv-store/app/util"
)

func main() {
	// Global setting
	timelibs.SetLocation(time.FixedZone("JST", 9*60*60))

	// Prepare grpc server settings
	sc, libErr := config.NewStorageConfig()
	if libErr != nil {
		util.FatalLog(fmt.Sprintf("Failed to load storage config: %v", libErr.Error()))
		return
	}
	storage, libErr := storage.NewStorage(sc)
	if libErr != nil {
		util.FatalLog(fmt.Sprintf("Failed to initialize storage: %v", libErr.Error()))
		return
	}

	server := service.NewGRPCServer(storage)
	lis, err := net.Listen("tcp", ":3000")
	if err != nil {
		util.FatalLog(fmt.Sprintf("Failed to listen: %v", err))
	}

	// Run
	go func() {
		if err := server.Serve(lis); err != nil {
			util.FatalLog("Failed to start gRPC server")
		}
	}()
	util.InfoLog("gRPC server started successfully")

	// Merge indexes
	go func() {
		ctx := context.Background()
		for {
			if libErr := storage.MergeIndexes(ctx); libErr != nil {
				util.ErrorLog(fmt.Sprintf("Failed to merge indexes: %v", libErr.Error()))
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// Shutdown settings
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
	<-quitCh
	// Stop accepting new request. This must be called before closeDB() method.
	server.GracefulStop()
	util.InfoLog("gRPC server stopped successfully")
}
