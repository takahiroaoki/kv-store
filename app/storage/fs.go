package storage

import (
	"os"
	"regexp"
	"sort"

	"github.com/takahiroaoki/go-libs/errorlibs"
)

const (
	logFilePrefix   = "log."
	indexFilePrefix = "index."
	csvExt          = ".csv"
	gobExt          = ".gob"
)

var (
	logFileNameRegexp = regexp.MustCompile(`log\.(\d+)\.csv`)
	logFilePathRegexp = regexp.MustCompile(`(.+)log\.(\d+)\.csv`)
)

func (s *storage) listFilesInAsc(dirPath string) (fileNameList []string, libErr errorlibs.Err) {
	fileNameList, libErr = s.listFiles(dirPath)
	if libErr != nil {
		return []string{}, nil
	}
	return sort.StringSlice(fileNameList), nil
}

func (s *storage) listFilesInDesc(dirPath string) (fileNameList []string, libErr errorlibs.Err) {
	fileNameList, libErr = s.listFiles(dirPath)
	if libErr != nil {
		return []string{}, nil
	}
	sort.Sort(sort.Reverse(sort.StringSlice(fileNameList)))
	return fileNameList, nil
}

func (s *storage) listFiles(dirPath string) (fileNameList []string, libErr errorlibs.Err) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}

	for _, f := range files {
		fileNameList = append(fileNameList, f.Name())
	}
	return fileNameList, nil
}

func (s *storage) createFile(filePath string) errorlibs.Err {
	f, err := os.Create(filePath)
	if err != nil {
		return errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	defer f.Close()
	return nil
}

func (s *storage) overwrite(oldFilePath, newFilePath string) errorlibs.Err {
	if err := os.Rename(oldFilePath, oldFilePath+".old"); err != nil {
		return errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	if err := os.Rename(newFilePath, oldFilePath); err != nil {
		return errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}
	_ = os.Remove(oldFilePath + ".old")
	return nil
}

func (s *storage) logFilePathToName(logFilePath string) string {
	return logFilePathRegexp.ReplaceAllString(logFilePath, logFilePrefix+"$2"+csvExt)
}
