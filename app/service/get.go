package service

import (
	"context"

	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func (s *kvServiceServer) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	kv, appErr := s.storage.GetByKey(ctx, req.GetKey())
	if appErr != nil {
		return nil, handleError(ctx, appErr)
	}
	return &pb.GetResponse{
		Value: kv.Value,
	}, nil
}
