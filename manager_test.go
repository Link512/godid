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
				PutFunc: func(e entry) error {
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

func TestGetToday(t *testing.T) {
	store = &entryStoreMock{
		GetRangeWithAggregationFunc: func(start, end time.Time, f aggregationFunction) (interface{}, error) {
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

func TestGetYesterday(t *testing.T) {
	store = &entryStoreMock{
		GetRangeWithAggregationFunc: func(start, end time.Time, f aggregationFunction) (interface{}, error) {
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
			GetRangeWithAggregationFunc: func(start, end time.Time, f aggregationFunction) (interface{}, error) {
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
			GetRangeWithAggregationFunc: func(start, end time.Time, f aggregationFunction) (interface{}, error) {
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
				GetRangeWithAggregationFunc: func(start, end time.Time, f aggregationFunction) (interface{}, error) {
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
