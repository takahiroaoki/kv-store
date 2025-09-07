package service

import (
	"github.com/takahiroaoki/kv-store/app/service/interceptor"
	"github.com/takahiroaoki/kv-store/app/storage"
	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

type kvServiceServer struct {
	pb.UnimplementedKVStoreServiceServer
	storage storage.Storage
}

func newKVServiceServer(storage storage.Storage) pb.KVStoreServiceServer {
	return &kvServiceServer{
		storage: storage,
	}
}

func NewGRPCServer(storage storage.Storage) *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(
		middleware.ChainUnaryServer(
			interceptor.SetContext(),
			interceptor.PerformanceLog(),
		),
	))
	reflection.Register(s)
	pb.RegisterKVStoreServiceServer(s, newKVServiceServer(storage))
	return s
}
