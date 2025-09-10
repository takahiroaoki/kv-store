package storage

import (
	"github.com/takahiroaoki/go-libs/errorlibs"
)

var (
	dataNotFound = errorlibs.NewErrFromMsg("data not found", errorlibs.CAUSE_NOT_FOUND, errorlibs.LOG_LEVEL_INFO)
)
