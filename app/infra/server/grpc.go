package server

import (
	"github.com/takahiroaoki/kv-store/app/infra/interceptor"
	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
)

type kvServiceServer struct {
	pb.UnimplementedKVStoreServiceServer
}

func newKVServiceServer() pb.KVStoreServiceServer {
	return &kvServiceServer{}
}

func NewGRPCServer() *grpc.Server {
	s := grpc.NewServer(grpc.UnaryInterceptor(
		middleware.ChainUnaryServer(
			interceptor.SetContext(),
			interceptor.PerformanceLog(),
		),
	))
	reflection.Register(s)
	pb.RegisterKVStoreServiceServer(s, newKVServiceServer())
	return s
}
