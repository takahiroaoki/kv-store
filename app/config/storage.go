package config

import "strings"

type StorageConfig struct {
	storageDir string
}

func (c *StorageConfig) StorageDir() string {
	return c.storageDir
}

func NewStorageConfig() StorageConfig {
	dir := env.storageDir
	if !strings.HasSuffix(env.storageDir, "/") {
		dir += "/"
	}
	return StorageConfig{
		storageDir: dir,
	}
}
