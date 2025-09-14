package config

import (
	"path/filepath"
	"strconv"

	"github.com/takahiroaoki/go-libs/errorlibs"
)

type StorageConfig struct {
	storageDir          string
	rowsPerLogFile      int
	maxPowerLogFile     int
	indexMergeBatchSize int
}

func (c *StorageConfig) StorageDir() string {
	return c.storageDir
}

func (c *StorageConfig) RowsPerLogFile() int {
	return c.rowsPerLogFile
}

func (c *StorageConfig) MaxPowerLogFile() int {
	return c.maxPowerLogFile
}

func (c *StorageConfig) IndexMergeBatchSize() int {
	return c.indexMergeBatchSize
}

func (c *StorageConfig) LogDir() string {
	return filepath.Join(c.storageDir, "logs")
}

func (c *StorageConfig) IndexDir() string {
	return filepath.Join(c.storageDir, "indexes")
}

func NewStorageConfig() (StorageConfig, errorlibs.Err) {
	rowsPerLogFile, err := strconv.Atoi(env.rowsPerLogFile)
	if err != nil {
		return StorageConfig{}, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	maxPowerLogFile, err := strconv.Atoi(env.maxPowerLogFile)
	if err != nil {
		return StorageConfig{}, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	indexMergeBatchSize, err := strconv.Atoi(env.indexMergeBatchSize)
	if err != nil {
		return StorageConfig{}, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	return StorageConfig{
		storageDir:          env.storageDir,
		rowsPerLogFile:      rowsPerLogFile,
		maxPowerLogFile:     maxPowerLogFile,
		indexMergeBatchSize: indexMergeBatchSize,
	}, nil
}
