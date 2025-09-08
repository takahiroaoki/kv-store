package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/takahiroaoki/kv-store/app/config"
	"github.com/takahiroaoki/kv-store/app/service"
	"github.com/takahiroaoki/kv-store/app/storage"
	"github.com/takahiroaoki/kv-store/app/util"
)

func main() {
	sc, appErr := config.NewStorageConfig()
	if appErr != nil {
		util.FatalLog(fmt.Sprintf("Failed to load storage config: %v", appErr.Error()))
		return
	}
	storage, appErr := storage.NewStorage(sc)
	if appErr != nil {
		util.FatalLog(fmt.Sprintf("Failed to initialize storage: %v", appErr.Error()))
		return
	}

	// Prepare grpc server settings
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

	// Shutdown settings
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGTERM, syscall.SIGINT)
	<-quitCh
	// Stop accepting new request. This must be called before closeDB() method.
	server.GracefulStop()
	util.InfoLog("gRPC server stopped successfully")
}
