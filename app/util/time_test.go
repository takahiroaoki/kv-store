package util

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_appTime_Format(t *testing.T) {
	t.Parallel()

	at := &appTime{
		time: time.Date(2023, 1, 2, 15, 4, 5, 123000000, time.UTC),
	}
	assert.Equal(t, "2023-01-02T15:04:05.123", at.Format())
}

func Test_appTime_Time(t *testing.T) {
	t.Parallel()

	at := &appTime{
		time: time.Date(2023, 1, 2, 15, 4, 5, 123000000, time.UTC),
	}
	assert.Equal(t, at.time, at.Time())
}
