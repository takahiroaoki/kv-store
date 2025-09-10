package storage

import (
	"os"
	"sort"

	"github.com/takahiroaoki/go-libs/errorlibs"
)

func (s *storage) listFilesInDesc(dirPath string) (fileNameList []string, libErr errorlibs.Err) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, errorlibs.NewErr(err, errorlibs.CAUSE_INTERNAL, errorlibs.LOG_LEVEL_ERROR)
	}

	for _, f := range files {
		fileNameList = append(fileNameList, f.Name())
	}
	sort.Sort(sort.Reverse(sort.StringSlice(fileNameList)))
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
