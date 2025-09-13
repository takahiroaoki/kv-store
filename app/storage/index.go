package storage

import (
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strconv"

	"github.com/takahiroaoki/go-libs/errorlibs"
	"github.com/takahiroaoki/kv-store/app/util"
)

type indexValue struct {
	FileName string
	Line     string // based on 1
}

type indexMap map[string]indexValue

func (s *storage) readIndex(ctx context.Context, indexFilePath string) (indexMap, errorlibs.Err) {
	idxMap := indexMap{}
	f, err := os.Open(indexFilePath)
	if err != nil {
		return idxMap, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	defer f.Close()
	dcd := gob.NewDecoder(f)

	if err := dcd.Decode(&idxMap); err != nil && err != io.EOF {
		return idxMap, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	return idxMap, nil
}

func (s *storage) updateIndex(ctx context.Context, key string, logFileName string, line int) errorlibs.Err {
	idxFilePath, libErr := s.indexFilePathFromLogFileName(logFileName)
	if libErr != nil {
		return libErr
	}

	idxMap, libErr := s.readIndex(ctx, idxFilePath)
	if libErr != nil {
		return libErr
	}
	idxMap[key] = indexValue{
		FileName: logFileName,
		Line:     strconv.Itoa(line),
	}
	util.InfoLog(fmt.Sprintf("new indexMap: %v", idxMap))

	tmpIdxFilePath := idxFilePath + ".tmp"
	if libErr := s.createFile(tmpIdxFilePath); libErr != nil {
		return libErr
	}
	tmpF, err := os.OpenFile(tmpIdxFilePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}

	ecd := gob.NewEncoder(tmpF)
	if err := ecd.Encode(idxMap); err != nil {
		return errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	defer tmpF.Close()

	if libErr := s.overwrite(idxFilePath, tmpIdxFilePath); libErr != nil {
		if err := os.Remove(tmpIdxFilePath); err != nil {
			util.WarnLog(err.Error())
		}
		return libErr
	}
	return nil
}

func (s *storage) indexFilePathFromLogFileName(logFileName string) (string, errorlibs.Err) {
	idxFileNames, libErr := s.listFiles(s.sc.IndexDir())
	if libErr != nil {
		return "", libErr
	}
	idxFileName := logFileNameRegexp.ReplaceAllString(logFileName, indexFilePrefix+"$1"+gobExt)
	idxFilePath := filepath.Join(s.sc.IndexDir(), idxFileName)
	if slices.Contains(idxFileNames, idxFileName) {
		return idxFilePath, nil
	}
	if libErr := s.createFile(idxFilePath); libErr != nil {
		return "", libErr
	}
	return idxFilePath, nil
}

func (s *storage) lookupLatestIndex(ctx context.Context, key string) (indexValue, errorlibs.Err) {
	idxFileNameList, libErr := s.listFilesInDesc(s.sc.IndexDir())
	if libErr != nil {
		return indexValue{}, libErr
	}
	if len(idxFileNameList) == 0 {
		return indexValue{}, dataNotFound
	}

	for _, idxFileName := range idxFileNameList {
		idxMap, libErr := s.readIndex(ctx, filepath.Join(s.sc.IndexDir(), idxFileName))
		if libErr != nil {
			return indexValue{}, libErr
		}
		if idxVal, ok := idxMap[key]; ok {
			return idxVal, nil
		}
	}
	return indexValue{}, dataNotFound
}
