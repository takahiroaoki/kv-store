package service

import (
	"context"

	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func (s *kvServiceServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if appErr := s.storage.DeleteKey(ctx, req.Key); appErr != nil {
		return nil, handleError(ctx, appErr)
	}
	return &pb.DeleteResponse{}, nil
}
