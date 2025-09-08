package config

import (
	"path/filepath"
	"strconv"

	"github.com/takahiroaoki/kv-store/app/util"
)

type StorageConfig struct {
	storageDir      string
	rowsPerLogFile  int
	maxPowerLogFile int
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

func (c *StorageConfig) LogDir() string {
	return filepath.Join(c.storageDir, "logs")
}

func NewStorageConfig() (StorageConfig, util.AppErr) {
	rowsPerLogFile, err := strconv.Atoi(env.rowsPerLogFile)
	if err != nil {
		return StorageConfig{}, util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	maxPowerLogFile, err := strconv.Atoi(env.maxPowerLogFile)
	if err != nil {
		return StorageConfig{}, util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	return StorageConfig{
		storageDir:      env.storageDir,
		rowsPerLogFile:  rowsPerLogFile,
		maxPowerLogFile: maxPowerLogFile,
	}, nil
}
