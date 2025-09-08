package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_PadStart(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		str    string
		length int
		padStr string
		want   string
	}{
		{
			name:   "Normal case",
			str:    "abc",
			length: 10,
			padStr: "ABC",
			want:   "ABCABCAabc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.want, PadStart(tt.str, tt.length, tt.padStr))
		})
	}
}
