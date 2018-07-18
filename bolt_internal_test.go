package godid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func timeFromString(t *testing.T, ts string) time.Time {
	result, err := time.Parse(time.RFC3339, ts)
	require.NoError(t, err)
	return result
}

func TestGetBucketRange(t *testing.T) {
	testCases := []struct {
		name        string
		start       time.Time
		end         time.Time
		expected    []string
		shouldError bool
	}{
		{
			name:        "empty",
			shouldError: true,
		},
		{
			name:        "bad interval",
			start:       time.Now().Add(1 * time.Hour),
			end:         time.Now(),
			shouldError: true,
		},
		{
			name:     "same day",
			start:    timeFromString(t, "2018-07-17T12:00:00Z"),
			end:      timeFromString(t, "2018-07-17T13:00:00Z"),
			expected: []string{"2018-07-17"},
		},
		{
			name:     "one day",
			start:    timeFromString(t, "2018-07-17T12:00:00Z"),
			end:      timeFromString(t, "2018-07-18T13:00:00Z"),
			expected: []string{"2018-07-17", "2018-07-18"},
		},
		{
			name:     "one week",
			start:    timeFromString(t, "2018-07-17T12:00:00Z"),
			end:      timeFromString(t, "2018-07-24T23:00:00Z"),
			expected: []string{"2018-07-17", "2018-07-18", "2018-07-19", "2018-07-20", "2018-07-21", "2018-07-22", "2018-07-23", "2018-07-24"},
		},
		{
			name:     "lower bound is limit",
			start:    timeFromString(t, "2018-07-17T00:00:00Z"),
			end:      timeFromString(t, "2018-07-18T23:00:00Z"),
			expected: []string{"2018-07-17", "2018-07-18"},
		},
		{
			name:     "upper bound is limit",
			start:    timeFromString(t, "2018-07-17T01:00:00Z"),
			end:      timeFromString(t, "2018-07-18T00:00:00Z"),
			expected: []string{"2018-07-17", "2018-07-18"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := getBucketRange(tc.start, tc.end)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}

func TestGetBucketFromTime(t *testing.T) {
	ts := timeFromString(t, "2018-07-18T00:00:00Z")
	bucket, err := getBucketFromTime(ts)
	require.NoError(t, err)
	assert.Equal(t, "2018-07-18", bucket)

	var empty time.Time
	_, err = getBucketFromTime(empty)
	require.Error(t, err)
}

func TestGetBucketFromEntry(t *testing.T) {
	entry := Entry{
		Timestamp: timeFromString(t, "2018-07-18T00:00:00Z"),
		Message:   []byte("asdb"),
	}
	bucket, err := getBucketFromEntry(entry)
	require.NoError(t, err)
	assert.Equal(t, "2018-07-18", bucket)

	empty := Entry{}
	_, err = getBucketFromEntry(empty)
	require.Error(t, err)
}
