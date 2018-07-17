package godid

import (
	"os"
	"testing"

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
	store, err := NewBoltStore("test.db")
	s.NoError(err)
	s.store = store.(*boltStore)
	s.db = s.store.db
}

func (s *boltTestSuite) TearDownTest() {
	s.store = nil
	os.Remove("test.db")
}

func (s *boltTestSuite) TestPut() {
	testCases := []struct {
		entries              []Entry
		shouldError          bool
		expectedBuckets      []string
		expectedBucketCounts map[string]int
	}{
		{
			entries:     []Entry{Entry{}},
			shouldError: true,
		},
		{
			entries: []Entry{
				Entry{
					Timestamp: timeFromString(s.T(), "2018-07-17T12:21:00Z"),
					Message:   []byte("BOOYAA"),
				},
			},
			expectedBuckets: []string{"2018-07-17"},
			expectedBucketCounts: map[string]int{
				"2018-07-17": 1,
			},
		},
		{
			entries: []Entry{
				Entry{
					Timestamp: timeFromString(s.T(), "2018-07-17T12:21:00Z"),
					Message:   []byte("BOOYAA"),
				},
				Entry{
					Timestamp: timeFromString(s.T(), "2018-07-18T12:21:00Z"),
					Message:   []byte("BOOYAAKAA"),
				},
			},
			expectedBuckets: []string{"2018-07-17", "2018-07-18"},
			expectedBucketCounts: map[string]int{
				"2018-07-17": 1,
				"2018-07-18": 1,
			},
		},
		{
			entries: []Entry{
				Entry{
					Timestamp: timeFromString(s.T(), "2018-07-17T12:21:00Z"),
					Message:   []byte("BOOYAA"),
				},
				Entry{
					Timestamp: timeFromString(s.T(), "2018-07-18T12:21:00Z"),
					Message:   []byte("BOOYAAKAA"),
				},
				Entry{
					Timestamp: timeFromString(s.T(), "2018-07-18T12:11:00Z"),
					Message:   []byte("FOOBAR"),
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

func TestBoltStore(t *testing.T) {
	suite.Run(t, new(boltTestSuite))
}
