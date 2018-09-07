package godid

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func TestAddEntry(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name: "empty",
		},
		{
			name:        "will error",
			shouldError: true,
		},
		{
			name:  "normal entry",
			input: "msg1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var insertedEntry entry
			store = &entryStoreMock{
				PutFunc: func(bucketName string, e entry) error {
					require.Equal(t, rootBucketName, bucketName)
					if tc.shouldError {
						return errors.New("BOOM")
					}
					insertedEntry = e
					return nil
				},
			}
			err := AddEntry(tc.input)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.False(t, insertedEntry.Timestamp.IsZero())
				require.Equal(t, []byte(tc.input), insertedEntry.Content)
			}
		})
	}
}

func TestAddEntryToBucket(t *testing.T) {
	testCases := []struct {
		name        string
		input       string
		bucketName  string
		shouldError bool
	}{
		{
			name: "empty",
		},
		{
			name:        "will error",
			shouldError: true,
		},
		{
			name:       "normal entry",
			bucketName: randString(10),
			input:      "msg1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var insertedEntry entry
			store = &entryStoreMock{
				PutFunc: func(bucketName string, e entry) error {
					require.Equal(t, tc.bucketName, bucketName)
					if tc.shouldError {
						return errors.New("BOOM")
					}
					insertedEntry = e
					return nil
				},
			}
			err := AddEntryToBucket(tc.bucketName, tc.input)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.False(t, insertedEntry.Timestamp.IsZero())
				require.Equal(t, []byte(tc.input), insertedEntry.Content)
			}
		})
	}
}

func TestGetToday(t *testing.T) {
	store = &entryStoreMock{
		GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
			require.Equal(t, rootBucketName, bucketName)
			curY, curM, curD := time.Now().Date()

			sY, sM, sD := start.Date()
			assert.Equal(t, curY, sY)
			assert.Equal(t, curM, sM)
			assert.Equal(t, curD, sD)

			eY, eM, eD := end.Date()
			assert.Equal(t, curY, eY)
			assert.Equal(t, curM, eM)
			assert.Equal(t, curD, eD)

			assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
			return []string{}, nil
		},
	}
	_, err := GetToday()
	require.NoError(t, err)
}

func TestGetTodayFromBucket(t *testing.T) {
	testBucketName := randString(10)
	store = &entryStoreMock{
		GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
			require.Equal(t, testBucketName, bucketName)
			curY, curM, curD := time.Now().Date()

			sY, sM, sD := start.Date()
			assert.Equal(t, curY, sY)
			assert.Equal(t, curM, sM)
			assert.Equal(t, curD, sD)

			eY, eM, eD := end.Date()
			assert.Equal(t, curY, eY)
			assert.Equal(t, curM, eM)
			assert.Equal(t, curD, eD)

			assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
			return []string{}, nil
		},
	}
	_, err := GetTodayFromBucket(testBucketName)
	require.NoError(t, err)
}

func TestGetYesterday(t *testing.T) {
	store = &entryStoreMock{
		GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
			require.Equal(t, rootBucketName, bucketName)
			curY, curM, curD := time.Now().AddDate(0, 0, -1).Date()

			sY, sM, sD := start.Date()
			assert.Equal(t, curY, sY)
			assert.Equal(t, curM, sM)
			assert.Equal(t, curD, sD)

			eY, eM, eD := end.Date()
			assert.Equal(t, curY, eY)
			assert.Equal(t, curM, eM)
			assert.Equal(t, curD, eD)

			assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
			return []string{}, nil
		},
	}
	_, err := GetYesterday()
	require.NoError(t, err)
}

func TestGetYesterdayFromBucket(t *testing.T) {
	testBucketName := randString(10)
	store = &entryStoreMock{
		GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
			require.Equal(t, testBucketName, bucketName)
			curY, curM, curD := time.Now().AddDate(0, 0, -1).Date()

			sY, sM, sD := start.Date()
			assert.Equal(t, curY, sY)
			assert.Equal(t, curM, sM)
			assert.Equal(t, curD, sD)

			eY, eM, eD := end.Date()
			assert.Equal(t, curY, eY)
			assert.Equal(t, curM, eM)
			assert.Equal(t, curD, eD)

			assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
			return []string{}, nil
		},
	}
	_, err := GetYesterdayFromBucket(testBucketName)
	require.NoError(t, err)
}

func TestGetThisWeek(t *testing.T) {
	testCases := []struct {
		flat bool
	}{
		{
			flat: true,
		},
		{
			flat: false,
		},
	}

	for _, tc := range testCases {
		store = &entryStoreMock{
			GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
				require.Equal(t, rootBucketName, bucketName)
				expectedStart, expectedEnd := getWeekInterval(time.Now()) //kinda circlejerking, but hey

				exY, exM, exD := expectedStart.Date()
				sY, sM, sD := start.Date()
				assert.Equal(t, exY, sY)
				assert.Equal(t, exM, sM)
				assert.Equal(t, exD, sD)

				exY, exM, exD = expectedEnd.Date()
				eY, eM, eD := end.Date()
				assert.Equal(t, exY, eY)
				assert.Equal(t, exM, eM)
				assert.Equal(t, exD, eD)

				if tc.flat {
					assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				} else {
					assert.Equal(t, reflect.ValueOf(perDayAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				}
				return nil, nil
			},
		}
		GetThisWeek(tc.flat)
	}
}

func TestGetThisWeekFromBucket(t *testing.T) {
	testCases := []struct {
		flat       bool
		bucketName string
	}{
		{
			flat:       true,
			bucketName: randString(10),
		},
		{
			flat:       false,
			bucketName: randString(10),
		},
	}

	for _, tc := range testCases {
		store = &entryStoreMock{
			GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
				require.Equal(t, tc.bucketName, bucketName)
				expectedStart, expectedEnd := getWeekInterval(time.Now()) //kinda circlejerking, but hey

				exY, exM, exD := expectedStart.Date()
				sY, sM, sD := start.Date()
				assert.Equal(t, exY, sY)
				assert.Equal(t, exM, sM)
				assert.Equal(t, exD, sD)

				exY, exM, exD = expectedEnd.Date()
				eY, eM, eD := end.Date()
				assert.Equal(t, exY, eY)
				assert.Equal(t, exM, eM)
				assert.Equal(t, exD, eD)

				if tc.flat {
					assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				} else {
					assert.Equal(t, reflect.ValueOf(perDayAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				}
				return nil, nil
			},
		}
		GetThisWeekFromBucket(tc.bucketName, tc.flat)
	}
}

func TestGetLastWeek(t *testing.T) {
	testCases := []struct {
		flat bool
	}{
		{
			flat: true,
		},
		{
			flat: false,
		},
	}

	for _, tc := range testCases {
		store = &entryStoreMock{
			GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
				require.Equal(t, rootBucketName, bucketName)
				expectedStart, expectedEnd := getWeekInterval(time.Now().AddDate(0, 0, -7)) //kinda circlejerking, but hey

				exY, exM, exD := expectedStart.Date()
				sY, sM, sD := start.Date()
				assert.Equal(t, exY, sY)
				assert.Equal(t, exM, sM)
				assert.Equal(t, exD, sD)

				exY, exM, exD = expectedEnd.Date()
				eY, eM, eD := end.Date()
				assert.Equal(t, exY, eY)
				assert.Equal(t, exM, eM)
				assert.Equal(t, exD, eD)

				if tc.flat {
					assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				} else {
					assert.Equal(t, reflect.ValueOf(perDayAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				}
				return nil, nil
			},
		}
		GetLastWeek(tc.flat)
	}
}

func TestGetLastWeekFromBucket(t *testing.T) {
	testCases := []struct {
		flat       bool
		bucketName string
	}{
		{
			flat:       true,
			bucketName: randString(10),
		},
		{
			flat:       false,
			bucketName: randString(10),
		},
	}

	for _, tc := range testCases {
		store = &entryStoreMock{
			GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
				require.Equal(t, tc.bucketName, bucketName)
				expectedStart, expectedEnd := getWeekInterval(time.Now().AddDate(0, 0, -7)) //kinda circlejerking, but hey

				exY, exM, exD := expectedStart.Date()
				sY, sM, sD := start.Date()
				assert.Equal(t, exY, sY)
				assert.Equal(t, exM, sM)
				assert.Equal(t, exD, sD)

				exY, exM, exD = expectedEnd.Date()
				eY, eM, eD := end.Date()
				assert.Equal(t, exY, eY)
				assert.Equal(t, exM, eM)
				assert.Equal(t, exD, eD)

				if tc.flat {
					assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				} else {
					assert.Equal(t, reflect.ValueOf(perDayAggregation).Pointer(), reflect.ValueOf(f).Pointer())
				}
				return nil, nil
			},
		}
		GetLastWeekFromBucket(tc.bucketName, tc.flat)
	}
}

func TestGetLastDuration(t *testing.T) {
	testCases := []struct {
		name        string
		interval    string
		flat        bool
		shouldError bool
	}{
		{
			name:        "bad interval",
			interval:    "1080h",
			shouldError: true,
		},
		{
			name:     "flat",
			interval: "45d",
			flat:     true,
		},
		{
			name:     "aggregated",
			interval: "45d",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store = &entryStoreMock{
				GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
					require.Equal(t, rootBucketName, bucketName)
					if !tc.shouldError {
						d, err := parseDuration(tc.interval)
						require.NoError(t, err)
						expectedStart := time.Now().Add(-1 * d)
						expectedEnd := time.Now()

						exY, exM, exD := expectedStart.Date()
						sY, sM, sD := start.Date()
						assert.Equal(t, exY, sY)
						assert.Equal(t, exM, sM)
						assert.Equal(t, exD, sD)

						exY, exM, exD = expectedEnd.Date()
						eY, eM, eD := end.Date()
						assert.Equal(t, exY, eY)
						assert.Equal(t, exM, eM)
						assert.Equal(t, exD, eD)

						if tc.flat {
							assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
						} else {
							assert.Equal(t, reflect.ValueOf(perDayAggregation).Pointer(), reflect.ValueOf(f).Pointer())
						}
					}
					if tc.flat {
						return []string{}, nil
					}
					return map[string][]string{}, nil
				},
			}
			_, err := GetLastDuration(tc.interval, tc.flat)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetLastDurationFromBucket(t *testing.T) {
	testCases := []struct {
		name        string
		interval    string
		bucketName  string
		flat        bool
		shouldError bool
	}{
		{
			name:        "bad interval",
			interval:    "1080h",
			shouldError: true,
		},
		{
			name:       "flat",
			interval:   "45d",
			bucketName: randString(10),
			flat:       true,
		},
		{
			name:       "aggregated",
			bucketName: randString(10),
			interval:   "45d",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			store = &entryStoreMock{
				GetRangeWithAggregationFunc: func(bucketName string, start, end time.Time, f aggregationFunction) (interface{}, error) {
					require.Equal(t, tc.bucketName, bucketName)
					if !tc.shouldError {
						d, err := parseDuration(tc.interval)
						require.NoError(t, err)
						expectedStart := time.Now().Add(-1 * d)
						expectedEnd := time.Now()

						exY, exM, exD := expectedStart.Date()
						sY, sM, sD := start.Date()
						assert.Equal(t, exY, sY)
						assert.Equal(t, exM, sM)
						assert.Equal(t, exD, sD)

						exY, exM, exD = expectedEnd.Date()
						eY, eM, eD := end.Date()
						assert.Equal(t, exY, eY)
						assert.Equal(t, exM, eM)
						assert.Equal(t, exD, eD)

						if tc.flat {
							assert.Equal(t, reflect.ValueOf(flatAggregation).Pointer(), reflect.ValueOf(f).Pointer())
						} else {
							assert.Equal(t, reflect.ValueOf(perDayAggregation).Pointer(), reflect.ValueOf(f).Pointer())
						}
					}
					if tc.flat {
						return []string{}, nil
					}
					return map[string][]string{}, nil
				},
			}
			_, err := GetLastDurationFromBucket(tc.bucketName, tc.interval, tc.flat)
			if tc.shouldError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
