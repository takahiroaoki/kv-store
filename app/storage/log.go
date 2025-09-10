package storage

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"

	"github.com/takahiroaoki/go-libs/stringlibs"
	"github.com/takahiroaoki/go-libs/timelibs"
	"github.com/takahiroaoki/kv-store/app/model"
	"github.com/takahiroaoki/kv-store/app/util"
)

type boolStr string

const (
	trueStr  boolStr = "1"
	falseStr boolStr = "0"
)

type logRow struct {
	key       string
	value     string
	delFlag   boolStr
	updatedAt string
}

func newLogRow(kv model.KeyValue, isDelete bool) logRow {
	delFlag := falseStr
	if isDelete {
		delFlag = trueStr
	}
	return logRow{
		key:       kv.Key,
		value:     kv.Value,
		delFlag:   delFlag,
		updatedAt: timelibs.Now().Format("2006-01-02T15:04:05"),
	}
}

func newKeyValueFromLogRow(row logRow) model.KeyValue {
	return model.KeyValue{
		Key:   row.key,
		Value: row.value,
	}
}

func (s *storage) nextLogFileName(currentLogFileName string) (string, util.AppErr) {
	prefix, postfix := "log.", ".csv"
	if len(currentLogFileName) == 0 {
		return prefix + stringlibs.PadStart("0", "0", s.sc.MaxPowerLogFile()) + postfix, nil
	}
	currNumStr := currentLogFileName[len(prefix) : len(currentLogFileName)-len(postfix)]
	currNum, err := strconv.Atoi(currNumStr)
	if err != nil {
		return "", util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	return prefix + stringlibs.PadStart(strconv.Itoa(currNum+1), "0", s.sc.MaxPowerLogFile()) + postfix, nil
}

func (s *storage) nextLogFilePath() (string, util.AppErr) {
	logFileNameList, appErr := s.listFilesInDesc(s.sc.LogDir())
	if appErr != nil {
		return "", appErr
	}
	if len(logFileNameList) == 0 {
		nextFileName, appErr := s.nextLogFileName("")
		if appErr != nil {
			return "", appErr
		}
		filePath := filepath.Join(s.sc.LogDir(), nextFileName)
		appErr = s.createFile(filePath)
		if appErr != nil {
			return "", appErr
		}
		return filePath, nil
	}

	latestFileName := logFileNameList[0]
	f, err := os.Open(filepath.Join(s.sc.LogDir(), latestFileName))
	if err != nil {
		return "", util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return "", util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}

	if len(records) < s.sc.RowsPerLogFile() {
		return filepath.Join(s.sc.LogDir(), latestFileName), nil
	}

	nextFileName, appErr := s.nextLogFileName(latestFileName)
	if appErr != nil {
		return "", appErr
	}
	nextFilePath := filepath.Join(s.sc.LogDir(), nextFileName)
	appErr = s.createFile(nextFilePath)
	if appErr != nil {
		return "", appErr
	}
	return nextFilePath, nil
}

func (s *storage) insertLogRow(ctx context.Context, row logRow) util.AppErr {
	path, appErr := s.nextLogFilePath()
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

	if err := writer.Write([]string{row.key, row.value, string(row.delFlag), row.updatedAt}); err != nil {
		return util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	return nil
}

func (s *storage) lookupTheLatestLogRow(ctx context.Context, key string) (logRow, util.AppErr) {
	logFileNameList, appErr := s.listFilesInDesc(s.sc.LogDir())
	if appErr != nil {
		return logRow{}, appErr
	}
	if len(logFileNameList) == 0 {
		return logRow{}, dataNotFound
	}

	for _, fileName := range logFileNameList {
		f, err := os.Open(filepath.Join(s.sc.LogDir(), fileName))
		if err != nil {
			return logRow{}, util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
		}
		defer f.Close()

		reader := csv.NewReader(f)
		records, err := reader.ReadAll()
		if err != nil {
			return logRow{}, util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
		}

		for i := len(records) - 1; i >= 0; i-- {
			if len(records[i]) < 4 {
				continue // skip illegal row
			}
			if records[i][0] == key {
				return logRow{
					key:       records[i][0],
					value:     records[i][1],
					delFlag:   boolStr(records[i][2]),
					updatedAt: records[i][3],
				}, nil
			}
		}
	}

	return logRow{}, dataNotFound
}

func (s *storage) InsertKeyValue(ctx context.Context, kv model.KeyValue) util.AppErr {
	return s.insertLogRow(ctx, newLogRow(kv, false))
}

func (s *storage) DeleteKey(ctx context.Context, key string) util.AppErr {
	return s.insertLogRow(ctx, newLogRow(model.KeyValue{Key: key}, true))
}

func (s *storage) GetByKey(ctx context.Context, key string) (model.KeyValue, util.AppErr) {
	logRow, appErr := s.lookupTheLatestLogRow(ctx, key)
	if appErr != nil {
		return model.KeyValue{}, appErr
	}
	if logRow.delFlag == trueStr {
		return model.KeyValue{}, dataNotFound
	}
	return newKeyValueFromLogRow(logRow), nil
}
