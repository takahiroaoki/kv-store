package service

import (
	"context"

	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func (s *kvServiceServer) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	if appErr := s.storage.InsertKeyValue(ctx, newKeyValue(req)); appErr != nil {
		return nil, handleError(ctx, appErr)
	}
	return &pb.SetResponse{}, nil
}
