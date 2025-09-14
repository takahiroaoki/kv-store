package config

import (
	"os"
)

var env *envVars

func init() {
	env = &envVars{
		storageDir:          getEnv("STORAGE_DIR", "/tmp/kv-store"), // Directory to store data. Default is /tmp/kv-store.
		rowsPerLogFile:      getEnv("ROWS_PER_LOG_FILE", "1000"),    // Number of rows per log file. Default is 1000.
		maxPowerLogFile:     getEnv("MAX_POWER_LOG_FILE", "10"),     // Log files can be created up to 10^n. Default is 10.
		indexMergeBatchSize: getEnv("INDEX_MERGE_BATCH_SIZE", "10"), // This should equals to or be larger than 2. Default is 10.
	}
}

type envVars struct {
	/* About FileSystem */
	storageDir          string
	rowsPerLogFile      string
	maxPowerLogFile     string
	indexMergeBatchSize string
}

func getEnv(name, defaultValue string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return defaultValue
}
