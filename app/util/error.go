package util

import (
	"errors"
	"fmt"
)

type ErrorCause int

const (
	CAUSE_UNDEFINED ErrorCause = iota
	CAUSE_INVALID_ARGUMENT
	CAUSE_NOT_FOUND
	CAUSE_INTERNAL
)

type LogLevel int

const (
	LOG_LEVEL_UNDEFINED LogLevel = iota
	LOG_LEVEL_INFO
	LOG_LEVEL_WARN
	LOG_LEVEL_ERROR
)

type AppErr interface {
	Error() string
	Cause() ErrorCause
	LogLevel() LogLevel
	AddErrContext(ctx string) AppErr
}

type appErr struct {
	err      error
	cause    ErrorCause
	logLevel LogLevel
}

func (de *appErr) Error() string {
	if de == nil || de.err == nil {
		return ""
	}
	return de.err.Error()
}

func (de *appErr) Cause() ErrorCause {
	if de == nil || de.err == nil {
		return CAUSE_UNDEFINED
	}
	return de.cause
}

func (de *appErr) LogLevel() LogLevel {
	if de == nil || de.err == nil {
		return LOG_LEVEL_UNDEFINED
	}
	return de.logLevel
}

func (de *appErr) AddErrContext(ctx string) AppErr {
	if de == nil {
		de = &appErr{}
	}
	de.err = fmt.Errorf("%s: %w", ctx, de)
	return de
}

func NewAppErr(err error, cause ErrorCause, logLevel LogLevel) AppErr {
	if err == nil {
		return nil
	}
	return &appErr{
		err:      err,
		cause:    cause,
		logLevel: logLevel,
	}
}

func NewAppErrFromMsg(msg string, cause ErrorCause, logLevel LogLevel) AppErr {
	return NewAppErr(errors.New(msg), cause, logLevel)
}
