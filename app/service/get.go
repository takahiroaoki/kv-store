package service

import (
	"context"

	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func (s *kvServiceServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	kv, libErr := s.storage.GetByKey(ctx, req.GetKey())
	if libErr != nil {
		return nil, handleError(ctx, libErr)
	}
	return &pb.GetResponse{
		Value: kv.Value,
	}, nil
}
