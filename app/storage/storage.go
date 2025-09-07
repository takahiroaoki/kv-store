package storage

import (
	"context"
	"os"

	"github.com/takahiroaoki/kv-store/app/config"
	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/util"
)

type Storage interface {
	InsertKeyValue(ctx context.Context, kv model.KeyValue) util.AppErr
}

type storage struct {
	sc config.StorageConfig
}

func NewStorage(sc config.StorageConfig) (Storage, util.AppErr) {
	if err := os.MkdirAll(sc.StorageDir(), 0755); err != nil {
		return nil, util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	return &storage{
		sc: sc,
	}, nil
}
