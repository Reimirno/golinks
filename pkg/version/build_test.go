package version

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildVersionString(t *testing.T) {
	testCases := []struct {
		name     string
		bv       BuildVersion
		expected string
	}{
		{
			name: "Complete BuildVersion",
			bv: BuildVersion{
				Version:   "1.0.0",
				Commit:    "abc123",
				BuildDate: "2023-04-15",
			},
			expected: fmt.Sprintf(
				versionFormat,
				"1.0.0",
				"abc123",
				"2023-04-15",
				runtime.Version(),
				runtime.GOOS,
				runtime.GOARCH,
			),
		},
		{
			name: "Empty BuildVersion",
			bv:   BuildVersion{},
			expected: fmt.Sprintf(
				versionFormat,
				"",
				"",
				"",
				runtime.Version(),
				runtime.GOOS,
				runtime.GOARCH,
			),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.bv.String()
			assert.Equal(t, tc.expected, result)
		})
	}
}
