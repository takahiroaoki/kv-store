package storage

import (
	"context"
	"os"
	"path/filepath"

	"github.com/takahiroaoki/go-libs/errorlibs"
	"github.com/takahiroaoki/go-libs/stringlibs"
	"github.com/takahiroaoki/kv-store/app/config"
	"github.com/takahiroaoki/kv-store/app/model"
)

type Storage interface {
	InsertKeyValue(ctx context.Context, kv model.KeyValue) errorlibs.Err
	DeleteKey(ctx context.Context, key string) errorlibs.Err
	GetByKey(ctx context.Context, key string) (model.KeyValue, errorlibs.Err)
	MergeIndexes(ctx context.Context) errorlibs.Err
}

type storage struct {
	sc config.StorageConfig
}

func (s *storage) setupFirstLogFile() errorlibs.Err {
	logFiles, libErr := s.listFiles(s.sc.LogDir())
	if libErr != nil {
		return libErr
	}

	if len(logFiles) > 0 {
		return nil
	}

	firstlogFilePath := filepath.Join(s.sc.LogDir(), logFilePrefix+stringlibs.PadStart("0", "0", s.sc.MaxPowerLogFile())+csvExt)
	if libErr := s.createFile(firstlogFilePath); libErr != nil {
		return libErr
	}
	return nil
}

func NewStorage(sc config.StorageConfig) (Storage, errorlibs.Err) {
	if err := os.MkdirAll(sc.StorageDir(), 0755); err != nil {
		return nil, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	if err := os.MkdirAll(sc.LogDir(), 0755); err != nil {
		return nil, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	if err := os.MkdirAll(sc.IndexDir(), 0755); err != nil {
		return nil, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	s := &storage{
		sc: sc,
	}
	if libErr := s.setupFirstLogFile(); libErr != nil {
		return nil, libErr
	}
	return s, nil
}
