package config

import (
	"os"
)

var env *envVars

func init() {
	env = &envVars{
		storageDir: os.Getenv("STORAGE_DIR"),
	}
}

type envVars struct {
	/* About FileSystem */
	storageDir string
}
