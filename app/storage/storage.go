package storage

import (
	"context"
	"os"

	"github.com/takahiroaoki/go-libs/errorlibs"
	"github.com/takahiroaoki/kv-store/app/config"
	"github.com/takahiroaoki/kv-store/app/model"
)

type Storage interface {
	InsertKeyValue(ctx context.Context, kv model.KeyValue) errorlibs.Err
	DeleteKey(ctx context.Context, key string) errorlibs.Err
	GetByKey(ctx context.Context, key string) (model.KeyValue, errorlibs.Err)
}

type storage struct {
	sc config.StorageConfig
}

func NewStorage(sc config.StorageConfig) (Storage, errorlibs.Err) {
	if err := os.MkdirAll(sc.StorageDir(), 0755); err != nil {
		return nil, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	if err := os.MkdirAll(sc.LogDir(), 0755); err != nil {
		return nil, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	return &storage{
		sc: sc,
	}, nil
}
