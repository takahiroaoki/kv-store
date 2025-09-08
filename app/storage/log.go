package storage

import (
	"context"
	"encoding/csv"
	"os"
	"time"

	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/util"
)

type logRow struct {
	key       string
	value     string
	delFlag   string
	updatedAt string
}

func newLogRowForSet(kv model.KeyValue) logRow {
	return logRow{
		key:       kv.Key,
		value:     kv.Value,
		delFlag:   "0",
		updatedAt: time.Now().Format("2006-01-02T15:04:05.123"),
	}
}

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

	row := newLogRowForSet(kv)
	if err := writer.Write([]string{row.key, row.value, row.delFlag, row.updatedAt}); err != nil {
		return util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	return nil
}
