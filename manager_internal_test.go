package godid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetPreviousWeekInterval(t *testing.T) {
	expectedStart := timeFromString(t, "2018-07-09T12:21:00Z")
	expectedEnd := timeFromString(t, "2018-07-15T12:21:00Z")
	reference := timeFromString(t, "2018-07-15T12:21:00Z")
	for i := 0; i < 7; i++ {
		reference = reference.AddDate(0, 0, 1)
		start, end := getPreviousWeekInterval(reference)
		assert.Equal(t, expectedStart, start)
		assert.Equal(t, expectedEnd, end)
	}
}

func TestGetRange(t *testing.T) {
	store = getTestBoltStore(t)
	testCases := []struct {
		name        string
		start       time.Time
		end         time.Time
		flat        bool
		shouldError bool
		expected    map[string][]string
	}{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := getRange(tc.start, tc.end, tc.flat)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
	cleanupTestBoltStore()
}
