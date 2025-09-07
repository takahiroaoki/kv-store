package service

import (
	"github.com/takahiroaoki/kv-store/app/model"
	pb "github.com/takahiroaoki/protobuf/gen_go/proto/kv_store/v1"
)

func newKeyValue(setReq *pb.SetRequest) model.KeyValue {
	return model.KeyValue{
		Key:   setReq.GetKey(),
		Value: setReq.GetValue(),
	}
}
