package util

import "time"

type AppTime interface {
	Format() string
	Time() time.Time
}

type appTime struct {
	time time.Time
}

func (t *appTime) Format() string {
	return t.time.Format("2006-01-02T15:04:05.123")
}

func (t *appTime) Time() time.Time {
	return t.time
}

func Since(t AppTime) time.Duration {
	return time.Since(t.Time())
}

var Now = func() AppTime {
	return &appTime{
		time: time.Now(),
	}
}
