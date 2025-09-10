package service

import (
	"context"

	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func (s *kvServiceServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	if libErr := s.storage.InsertKeyValue(ctx, newKeyValue(req)); libErr != nil {
		return nil, handleError(ctx, libErr)
	}
	return &pb.SetResponse{}, nil
}
