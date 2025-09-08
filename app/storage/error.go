package storage

import (
	"errors"

	"github.com/takahiroaoki/kv-store/app/util"
)

var (
	dataNotFound = util.NewAppErr(errors.New("data not found"), util.CAUSE_NOT_FOUND, util.LOG_LEVEL_INFO)
)
