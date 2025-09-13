package storage

import (
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"

	"github.com/takahiroaoki/go-libs/errorlibs"
	"github.com/takahiroaoki/go-libs/stringlibs"
	"github.com/takahiroaoki/go-libs/timelibs"
	"github.com/takahiroaoki/kv-store/app/model"
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

func (s *storage) nextLogFileName(currentLogFileName string) (string, errorlibs.Err) {
	currNumStr := currentLogFileName[len(logFilePrefix) : len(currentLogFileName)-len(csvExt)]
	currNum, err := strconv.Atoi(currNumStr)
	if err != nil {
		return "", errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	return logFilePrefix + stringlibs.PadStart(strconv.Itoa(currNum+1), "0", s.sc.MaxPowerLogFile()) + csvExt, nil
}

func (s *storage) nextLogFilePath() (string, int, errorlibs.Err) { // (logFilePath string, len(records) int, libErr errorlibs.Err)
	logFileNameList, libErr := s.listFilesInDesc(s.sc.LogDir())
	if libErr != nil {
		return "", 0, libErr
	}
	if len(logFileNameList) == 0 {
		return "", 0, errorlibs.NewErrFromMsg("not setup yet", errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}

	latestFileName := logFileNameList[0]
	f, err := os.Open(filepath.Join(s.sc.LogDir(), latestFileName))
	if err != nil {
		return "", 0, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	if err != nil {
		return "", 0, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}

	if len(records) < s.sc.RowsPerLogFile() {
		return filepath.Join(s.sc.LogDir(), latestFileName), len(records), nil
	}

	nextFileName, libErr := s.nextLogFileName(latestFileName)
	if libErr != nil {
		return "", 0, libErr
	}
	nextFilePath := filepath.Join(s.sc.LogDir(), nextFileName)
	libErr = s.createFile(nextFilePath)
	if libErr != nil {
		return "", 0, libErr
	}
	return nextFilePath, 0, nil
}

func (s *storage) insertLogRow(ctx context.Context, row logRow) errorlibs.Err {
	logFilePath, records, libErr := s.nextLogFilePath()
	if libErr != nil {
		return libErr
	}
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	defer logFile.Close()

	writer := csv.NewWriter(logFile)
	defer writer.Flush()

	if err := writer.Write([]string{row.key, row.value, string(row.delFlag), row.updatedAt}); err != nil {
		return errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}

	if libErr := s.updateIndex(row.key, s.logFilePathToName(logFilePath), records+1); libErr != nil {
		return libErr
	}
	return nil
}

func (s *storage) lookupTheLatestLogRow(ctx context.Context, key string) (logRow, errorlibs.Err) {
	logFileNameList, libErr := s.listFilesInDesc(s.sc.LogDir())
	if libErr != nil {
		return logRow{}, libErr
	}
	if len(logFileNameList) == 0 {
		return logRow{}, dataNotFound
	}

	for _, fileName := range logFileNameList {
		f, err := os.Open(filepath.Join(s.sc.LogDir(), fileName))
		if err != nil {
			return logRow{}, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
		}
		defer f.Close()

		reader := csv.NewReader(f)
		records, err := reader.ReadAll()
		if err != nil {
			return logRow{}, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
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

func (s *storage) InsertKeyValue(ctx context.Context, kv model.KeyValue) errorlibs.Err {
	return s.insertLogRow(ctx, newLogRow(kv, false))
}

func (s *storage) DeleteKey(ctx context.Context, key string) errorlibs.Err {
	return s.insertLogRow(ctx, newLogRow(model.KeyValue{Key: key}, true))
}

func (s *storage) GetByKey(ctx context.Context, key string) (model.KeyValue, errorlibs.Err) {
	logRow, libErr := s.lookupTheLatestLogRow(ctx, key)
	if libErr != nil {
		return model.KeyValue{}, libErr
	}
	if logRow.delFlag == trueStr {
		return model.KeyValue{}, dataNotFound
	}
	return newKeyValueFromLogRow(logRow), nil
}
