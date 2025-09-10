package service

import (
	"context"

	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func (s *kvServiceServer) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	if libErr := s.storage.DeleteKey(ctx, req.Key); libErr != nil {
		return nil, handleError(ctx, libErr)
	}
	return &pb.DeleteResponse{}, nil
}
