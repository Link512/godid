package godid

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWeekInterval(t *testing.T) {
	expectedStart := timeFromString(t, "2018-07-09T12:21:00Z")
	expectedEnd := timeFromString(t, "2018-07-15T12:21:00Z")
	reference := timeFromString(t, "2018-07-08T12:21:00Z")
	for i := 0; i < 7; i++ {
		reference = reference.AddDate(0, 0, 1)
		start, end := getWeekInterval(reference)
		assert.Equal(t, expectedStart, start)
		assert.Equal(t, expectedEnd, end)
	}
}

func TestGetRange(t *testing.T) {
	testCases := []struct {
		name             string
		start            time.Time
		end              time.Time
		flat             bool
		storeReturn      interface{}
		storeShouldError bool
		shouldError      bool
		expected         map[string][]string
	}{
		{
			name:             "store error",
			start:            time.Now(),
			end:              time.Now().AddDate(0, 0, 1),
			shouldError:      true,
			storeShouldError: true,
		},
		{
			name:        "bad flat result from store",
			start:       time.Now(),
			end:         time.Now().AddDate(0, 0, 1),
			flat:        true,
			storeReturn: []int{1, 2, 3},
			shouldError: true,
		},
		{
			name:        "bad aggregated result from store",
			start:       time.Now(),
			end:         time.Now().AddDate(0, 0, 1),
			storeReturn: []int{1, 2, 3},
			shouldError: true,
		},
		{
			name:        "flat",
			start:       time.Now(),
			end:         time.Now().AddDate(0, 0, 1),
			flat:        true,
			storeReturn: []string{"a", "b", "c"},
			expected:    map[string][]string{"flat": []string{"a", "b", "c"}},
		},
		{
			name:  "per day",
			start: time.Now(),
			end:   time.Now().AddDate(0, 0, 1),
			storeReturn: map[string][]string{
				"key1": []string{"a", "b", "c"},
				"key2": []string{"foo"},
			},
			expected: map[string][]string{
				"key1": []string{"a", "b", "c"},
				"key2": []string{"foo"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store = &entryStoreMock{
				GetRangeWithAggregationFunc: func(_, _ time.Time, f aggregationFunction) (interface{}, error) {
					if tc.storeShouldError {
						return nil, errors.New("boom")
					}
					var expected uintptr
					if tc.flat {
						expected = reflect.ValueOf(flatAggregation).Pointer()
					} else {
						expected = reflect.ValueOf(perDayAggregation).Pointer()
					}
					actual := reflect.ValueOf(f).Pointer()
					require.Equal(t, expected, actual)
					return tc.storeReturn, nil
				},
			}
			result, err := getRange(tc.start, tc.end, tc.flat)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expected, result)
			}
		})
	}
}

func TestFlatAggregation(t *testing.T) {
	testCases := []struct {
		name     string
		input    []entry
		expected []string
	}{
		{
			name:     "empty",
			expected: []string{},
		},
		{
			name: "entries in one day",
			input: []entry{
				entry{
					Timestamp: timeFromString(t, "2018-07-18T12:00:00Z"),
					Content:   []byte("msg1"),
				},
				entry{
					Timestamp: timeFromString(t, "2018-07-18T14:00:00Z"),
					Content:   []byte("msg2"),
				},
			},
			expected: []string{"msg1", "msg2"},
		},
		{
			name: "entries in separate days",
			input: []entry{
				entry{
					Timestamp: timeFromString(t, "2018-07-18T12:00:00Z"),
					Content:   []byte("msg1"),
				},
				entry{
					Timestamp: timeFromString(t, "2018-07-18T14:00:00Z"),
					Content:   []byte("msg2"),
				},
				entry{
					Timestamp: timeFromString(t, "2018-07-20T14:00:00Z"),
					Content:   []byte("msg3"),
				},
			},
			expected: []string{"msg1", "msg2", "msg3"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := flatAggregation(tc.input)
			require.NoError(t, err)
			actual, ok := result.([]string)
			require.True(t, ok)
			require.Equal(t, tc.expected, actual)

		})
	}
}

func TestPerDayAggregation(t *testing.T) {
	testCases := []struct {
		name        string
		input       []entry
		shouldError bool
		expected    map[string][]string
	}{
		{
			name:     "empty",
			expected: map[string][]string{},
		},
		{
			name: "bad entry",
			input: []entry{
				entry{
					Content: []byte("bad"),
				},
			},
			shouldError: true,
		},
		{
			name: "entries in one day",
			input: []entry{
				entry{
					Timestamp: timeFromString(t, "2018-07-18T12:00:00Z"),
					Content:   []byte("msg1"),
				},
				entry{
					Timestamp: timeFromString(t, "2018-07-18T14:00:00Z"),
					Content:   []byte("msg2"),
				},
			},
			expected: map[string][]string{
				"2018-07-18": []string{"msg1", "msg2"},
			},
		},
		{
			name: "entries in separate days",
			input: []entry{
				entry{
					Timestamp: timeFromString(t, "2018-07-18T12:00:00Z"),
					Content:   []byte("msg1"),
				},
				entry{
					Timestamp: timeFromString(t, "2018-07-18T14:00:00Z"),
					Content:   []byte("msg2"),
				},
				entry{
					Timestamp: timeFromString(t, "2018-07-20T14:00:00Z"),
					Content:   []byte("msg3"),
				},
			},
			expected: map[string][]string{
				"2018-07-18": []string{"msg1", "msg2"},
				"2018-07-20": []string{"msg3"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := perDayAggregation(tc.input)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				actual, ok := result.(map[string][]string)
				require.True(t, ok)
				require.Equal(t, tc.expected, actual)
			}
		})
	}
}
