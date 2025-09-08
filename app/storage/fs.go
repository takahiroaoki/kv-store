package storage

import (
	"os"
	"sort"

	"github.com/takahiroaoki/kv-store/app/util"
)

func (s *storage) listFilesInDesc(dirPath string) (fileNameList []string, appErr util.AppErr) {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}

	for _, f := range files {
		fileNameList = append(fileNameList, f.Name())
	}
	sort.Sort(sort.Reverse(sort.StringSlice(fileNameList)))
	return fileNameList, nil
}

func (s *storage) createFile(filePath string) util.AppErr {
	f, err := os.Create(filePath)
	if err != nil {
		return util.NewAppErr(err, util.CAUSE_INTERNAL, util.LOG_LEVEL_ERROR)
	}
	defer f.Close()
	return nil
}
