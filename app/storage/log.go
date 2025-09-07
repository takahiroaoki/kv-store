package storage

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/util"
)

func (s *storage) currentLogFilePath() (string, util.AppErr) {
	fileName := "log.csv"
	if _, err := os.Stat(s.sc.StorageDir() + fileName); os.IsNotExist(err) {
		f, err := os.Create(s.sc.StorageDir() + fileName)
		if err != nil {
			return "", util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
		}
		defer f.Close()
	}

	return s.sc.StorageDir() + fileName, nil
}

func (s *storage) InsertKeyValue(ctx context.Context, kv model.KeyValue) util.AppErr {
	path, appErr := s.currentLogFilePath()
	if appErr != nil {
		return appErr
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	if err := writer.Write([]string{kv.Key, kv.Value, "0"}); err != nil {
		return util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	return nil
}
