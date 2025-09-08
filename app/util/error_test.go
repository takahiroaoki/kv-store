package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_appErr_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ae       *appErr
		expected string
	}{
		{
			name: "Success",
			ae: &appErr{
				err: errors.New("err"),
			},
			expected: "err",
		},
		{
			name:     "Success(*appErr is nil)",
			ae:       nil,
			expected: "",
		},
		{
			name:     "Success(err is nil)",
			ae:       &appErr{},
			expected: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.ae.Error())
		})
	}
}

func Test_appErr_Cause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ae       *appErr
		expected ErrorCause
	}{
		{
			name: "Success",
			ae: &appErr{
				err:   errors.New("err"),
				cause: CAUSE_INTERNAL,
			},
			expected: CAUSE_INTERNAL,
		},
		{
			name:     "Success(*appErr is nil)",
			ae:       nil,
			expected: CAUSE_UNDEFINED,
		},
		{
			name:     "Success(err is nil)",
			ae:       &appErr{},
			expected: CAUSE_UNDEFINED,
		},
		{
			name: "Error(cause is not defined)",
			ae: &appErr{
				err: errors.New("err"),
			},
			expected: CAUSE_UNDEFINED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.ae.Cause())
		})
	}
}

func Test_appErr_LogLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		ae       *appErr
		expected LogLevel
	}{
		{
			name: "Success",
			ae: &appErr{
				err:      errors.New("err"),
				logLevel: LOG_LEVEL_INFO,
			},
			expected: LOG_LEVEL_INFO,
		},
		{
			name:     "Success(*appErr is nil)",
			ae:       nil,
			expected: LOG_LEVEL_UNDEFINED,
		},
		{
			name:     "Success(err is nil)",
			ae:       &appErr{},
			expected: LOG_LEVEL_UNDEFINED,
		},
		{
			name: "Error(logLevel is not defined)",
			ae: &appErr{
				err: errors.New("err"),
			},
			expected: LOG_LEVEL_UNDEFINED,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expected, tt.ae.LogLevel())
		})
	}
}

func Test_appErr_AddErrContext(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx string
	}

	tests := []struct {
		name     string
		ae       *appErr
		args     args
		expected *appErr
	}{
		{
			name: "Success",
			ae: &appErr{
				err:      errors.New("err"),
				logLevel: LOG_LEVEL_INFO,
			},
			args: args{
				ctx: "ctx",
			},
			expected: &appErr{
				err:      errors.New("ctx: err"),
				logLevel: LOG_LEVEL_INFO,
			},
		},
		{
			name: "Success(*appErr is nil)",
			ae:   nil,
			args: args{
				ctx: "ctx",
			},
			expected: &appErr{
				err: errors.New("ctx: "),
			},
		},
		{
			name: "Success(err is nil)",
			ae:   &appErr{},
			args: args{
				ctx: "ctx",
			},
			expected: &appErr{
				err: errors.New("ctx: "),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.ae.AddErrContext(tt.args.ctx)
			assert.EqualError(t, tt.expected, got.Error())
		})
	}
}
