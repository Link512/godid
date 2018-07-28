package godid

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/boltdb/bolt"

	"github.com/stretchr/testify/suite"
)

type boltTestSuite struct {
	suite.Suite
	store *boltStore
	db    *bolt.DB
}

func (s *boltTestSuite) SetupSuite() {
	os.Remove("test.db")
}

func (s *boltTestSuite) SetupTest() {
	s.store = getTestBoltStore(s.T())
	s.db = s.store.db
}

func (s *boltTestSuite) TearDownTest() {
	s.store.Close()
	s.store = nil
	os.Remove("test.db")
}

func (s *boltTestSuite) TestPut() {
	testCases := []struct {
		entries              []entry
		shouldError          bool
		expectedBuckets      []string
		expectedBucketCounts map[string]int
	}{
		{
			entries:     []entry{entry{}},
			shouldError: true,
		},
		{
			entries: []entry{
				entry{
					Timestamp: timeFromString(s.T(), "2018-07-17T12:21:00Z"),
					Content:   []byte("BOOYAA"),
				},
			},
			expectedBuckets: []string{"2018-07-17"},
			expectedBucketCounts: map[string]int{
				"2018-07-17": 1,
			},
		},
		{
			entries: []entry{
				entry{
					Timestamp: timeFromString(s.T(), "2018-07-17T12:21:00Z"),
					Content:   []byte("BOOYAA"),
				},
				entry{
					Timestamp: timeFromString(s.T(), "2018-07-18T12:21:00Z"),
					Content:   []byte("BOOYAAKAA"),
				},
			},
			expectedBuckets: []string{"2018-07-17", "2018-07-18"},
			expectedBucketCounts: map[string]int{
				"2018-07-17": 1,
				"2018-07-18": 1,
			},
		},
		{
			entries: []entry{
				entry{
					Timestamp: timeFromString(s.T(), "2018-07-17T12:21:00Z"),
					Content:   []byte("BOOYAA"),
				},
				entry{
					Timestamp: timeFromString(s.T(), "2018-07-18T12:21:00Z"),
					Content:   []byte("BOOYAAKAA"),
				},
				entry{
					Timestamp: timeFromString(s.T(), "2018-07-18T12:11:00Z"),
					Content:   []byte("FOOBAR"),
				},
			},
			expectedBuckets: []string{"2018-07-17", "2018-07-18"},
			expectedBucketCounts: map[string]int{
				"2018-07-17": 1,
				"2018-07-18": 2,
			},
		},
	}

	for _, tc := range testCases {
		for _, entry := range tc.entries {
			err := s.store.Put(entry)
			if tc.shouldError {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		}
		s.db.View(func(tx *bolt.Tx) error {
			for _, b := range tc.expectedBuckets {
				bucket := tx.Bucket([]byte(b))
				s.NotNil(bucket)
				s.Equal(tc.expectedBucketCounts[b], bucket.Stats().KeyN)
			}
			return nil
		})
	}
}

func (s *boltTestSuite) TestGetRange() {
	entries := []entry{
		entry{Timestamp: timeFromString(s.T(), "2018-07-18T12:11:00Z"), Content: []byte("msg1")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-18T13:32:00Z"), Content: []byte("msg2")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-18T14:21:00Z"), Content: []byte("msg3")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-19T09:11:00Z"), Content: []byte("msg4")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-19T23:15:00Z"), Content: []byte("msg5")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-20T10:11:00Z"), Content: []byte("msg6")},
	}
	for _, entry := range entries {
		err := s.store.Put(entry)
		s.NoError(err)
	}
	testCases := []struct {
		start       time.Time
		end         time.Time
		shouldError bool
		expected    []entry
	}{
		{
			shouldError: true,
		},
		{
			start:       time.Now().Add(1 * time.Hour),
			end:         time.Now(),
			shouldError: true,
		},
		{
			start:    timeFromString(s.T(), "2018-06-20T10:11:00Z"),
			end:      timeFromString(s.T(), "2018-06-21T10:11:00Z"),
			expected: []entry{},
		},
		{
			start:    timeFromString(s.T(), "2018-07-20T09:11:00Z"),
			end:      timeFromString(s.T(), "2018-07-20T10:11:00Z"),
			expected: []entry{entry{Timestamp: timeFromString(s.T(), "2018-07-20T10:11:00Z"), Content: []byte("msg6")}},
		},
		{
			start: timeFromString(s.T(), "2018-07-18T12:11:00Z"),
			end:   timeFromString(s.T(), "2018-07-19T09:11:00Z"),
			expected: []entry{
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T12:11:00Z"), Content: []byte("msg1")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T13:32:00Z"), Content: []byte("msg2")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T14:21:00Z"), Content: []byte("msg3")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-19T09:11:00Z"), Content: []byte("msg4")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-19T23:15:00Z"), Content: []byte("msg5")},
			},
		},
		{
			start: timeFromString(s.T(), "2018-06-18T12:11:00Z"),
			end:   timeFromString(s.T(), "2018-09-19T09:11:00Z"),
			expected: []entry{
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T12:11:00Z"), Content: []byte("msg1")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T13:32:00Z"), Content: []byte("msg2")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T14:21:00Z"), Content: []byte("msg3")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-19T09:11:00Z"), Content: []byte("msg4")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-19T23:15:00Z"), Content: []byte("msg5")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-20T10:11:00Z"), Content: []byte("msg6")},
			},
		},
	}

	for _, tc := range testCases {
		entries, err := s.store.GetRange(tc.start, tc.end)
		if tc.shouldError {
			s.Error(err)
		} else {
			s.NoError(err)
			s.Equal(tc.expected, entries)
		}
	}
}

func (s *boltTestSuite) TestGetRangeWithAggregation() {
	entries := []entry{
		entry{Timestamp: timeFromString(s.T(), "2018-07-18T12:11:00Z"), Content: []byte("msg1")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-18T13:32:00Z"), Content: []byte("msg2")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-18T14:21:00Z"), Content: []byte("msg3")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-19T09:11:00Z"), Content: []byte("msg4")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-19T23:15:00Z"), Content: []byte("msg5")},
		entry{Timestamp: timeFromString(s.T(), "2018-07-20T10:11:00Z"), Content: []byte("msg6")},
	}
	for _, entry := range entries {
		err := s.store.Put(entry)
		s.NoError(err)
	}
	testCases := []struct {
		start       time.Time
		end         time.Time
		agg         aggregationFunction
		shouldError bool
		expected    interface{}
	}{
		{
			shouldError: true,
		},
		{
			start:       time.Now().Add(1 * time.Hour),
			end:         time.Now(),
			shouldError: true,
		},
		{
			start:       timeFromString(s.T(), "2018-06-20T10:11:00Z"),
			end:         timeFromString(s.T(), "2018-06-21T10:11:00Z"),
			shouldError: true,
		},
		{
			start: timeFromString(s.T(), "2018-06-20T10:11:00Z"),
			end:   timeFromString(s.T(), "2018-06-21T10:11:00Z"),
			agg: func(e []entry) (interface{}, error) {
				return nil, errors.New("BOOM")
			},
			shouldError: true,
		},
		{
			start: timeFromString(s.T(), "2018-06-18T12:11:00Z"),
			end:   timeFromString(s.T(), "2018-09-19T09:11:00Z"),
			agg: func(e []entry) (interface{}, error) {
				return e, nil
			},
			expected: []entry{
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T12:11:00Z"), Content: []byte("msg1")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T13:32:00Z"), Content: []byte("msg2")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-18T14:21:00Z"), Content: []byte("msg3")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-19T09:11:00Z"), Content: []byte("msg4")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-19T23:15:00Z"), Content: []byte("msg5")},
				entry{Timestamp: timeFromString(s.T(), "2018-07-20T10:11:00Z"), Content: []byte("msg6")},
			},
		},
	}

	for _, tc := range testCases {
		result, err := s.store.GetRangeWithAggregation(tc.start, tc.end, tc.agg)
		if tc.shouldError {
			s.Error(err)
		} else {
			s.NoError(err)
			entries, ok := result.([]entry)
			s.True(ok)
			s.Equal(tc.expected, entries)
		}
	}
}

func TestBoltStore(t *testing.T) {
	suite.Run(t, new(boltTestSuite))
}
