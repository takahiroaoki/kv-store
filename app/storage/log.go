package storage

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/util"
)

type logRow struct {
	key       string
	value     string
	delFlag   string
	updatedAt string
}

func newLogRow(kv model.KeyValue, isDelete bool) logRow {
	delFlag := "0"
	if isDelete {
		delFlag = "1"
	}
	return logRow{
		key:       kv.Key,
		value:     kv.Value,
		delFlag:   delFlag,
		updatedAt: util.Now().Format(),
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

func (s *storage) insertLogRow(ctx context.Context, row logRow) util.AppErr {
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

	if err := writer.Write([]string{row.key, row.value, row.delFlag, row.updatedAt}); err != nil {
		return util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	return nil
}

func (s *storage) InsertKeyValue(ctx context.Context, kv model.KeyValue) util.AppErr {
	return s.insertLogRow(ctx, newLogRow(kv, false))
}

func (s *storage) DeleteKey(ctx context.Context, key string) util.AppErr {
	return s.insertLogRow(ctx, newLogRow(model.KeyValue{Key: key}, true))
}
