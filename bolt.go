package godid

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
)

type boltStore struct {
	db *bolt.DB
}

const (
	timeFormat = time.RFC3339
)

//NewBoltStore creates new EntryStore with boltdb as a backend
//TODO: don't export this, use NewStore()
func NewBoltStore(fileName string) (EntryStore, error) {
	db, err := bolt.Open(fileName, 0600, nil)
	if err != nil {
		return nil, err
	}
	return &boltStore{
		db: db,
	}, nil
}

func (s *boltStore) Put(e Entry) error {
	bucketName, err := getBucketFromEntry(e)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(e.Timestamp.Format(timeFormat)), e.Message)
	})
}

func (s *boltStore) GetRange(start, end time.Time) ([]Entry, error) {
	buckets, err := getBucketRange(start, end)
	if err != nil {
		return nil, err
	}
	result := make([]Entry, 0)
	for _, bucket := range buckets {
		err := s.db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(bucket))
			if b == nil {
				return nil
			}
			return b.ForEach(func(k, v []byte) error {
				timestamp, err := time.Parse(timeFormat, string(k))
				if err != nil {
					return err
				}
				result = append(result, Entry{
					Timestamp: timestamp,
					Message:   v,
				})
				return nil
			})
		})
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (s *boltStore) GetRangeWithAggregation(start, end time.Time, agg AggregationFunction) (interface{}, error) {
	if agg == nil {
		return nil, errors.New("aggregation function is nil")
	}
	entries, err := s.GetRange(start, end)
	if err != nil {
		return nil, err
	}
	return agg(entries)
}

func getBucketFromEntry(e Entry) (string, error) {
	return getBucketFromTime(e.Timestamp)
}

func getBucketFromTime(t time.Time) (string, error) {
	if t.IsZero() {
		return "", errors.New("timestamp can't be zero")
	}
	return t.Format("2006-01-02"), nil
}

func getBucketRange(start, end time.Time) ([]string, error) {
	if start.IsZero() || end.IsZero() {
		return nil, errors.New("start and end must be set")
	}
	if start.After(end) {
		return nil, errors.New("start time is after end time")
	}

	strippedStart := start.Truncate(24 * time.Hour)
	strippedEnd := end.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	buckets := make([]string, 0)
	for i := strippedStart.UnixNano(); i < strippedEnd.UnixNano(); i += int64(24 * time.Hour) {
		bucket, err := getBucketFromTime(time.Unix(0, i))
		if err != nil {
			return nil, err
		}
		buckets = append(buckets, bucket)
	}
	return buckets, nil
}
