package file_system

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size     int64
		expected string
	}{
		{1023, "1023 B"},
		{2048, "2.00 KB"},
		{3145728, "3.00 MB"},
		{5368709120, "5.00 GB"},
		{1099511627776, "1.00 TB"},
		{1125899906842624, "1.00 PB"},
	}

	for _, test := range tests {
		result := formatFileSize(test.size)
		assert.Equal(t, test.expected, result)
	}
}
