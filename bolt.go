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

// newBoltStore creates new entryStore with boltdb as a backend
func newBoltStore(cfg config) (*boltStore, error) {
	path, err := cfg.GetStorePath()
	if err != nil {
		return nil, err
	}
	db, err := bolt.Open(path, 0600, &bolt.Options{Timeout: 10 * time.Second})
	if err != nil {
		return nil, err
	}
	return &boltStore{
		db: db,
	}, nil
}

func (s *boltStore) Put(parentBucketName string, e entry) error {
	bucketName, err := getBucketFromEntry(e)
	if err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		parentBucket, err := tx.CreateBucketIfNotExists([]byte(parentBucketName))
		if err != nil {
			return err
		}
		b, err := parentBucket.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return err
		}
		return b.Put([]byte(e.Timestamp.Format(timeFormat)), e.Content)
	})
}

func (s *boltStore) GetRange(parentBucketName string, start, end time.Time) ([]entry, error) {
	buckets, err := getBucketRange(start, end)
	if err != nil {
		return nil, err
	}
	result := make([]entry, 0)
	for _, bucket := range buckets {
		err := s.db.View(func(tx *bolt.Tx) error {
			parentBucket := tx.Bucket([]byte(parentBucketName))
			if parentBucket == nil {
				return nil
			}
			b := parentBucket.Bucket([]byte(bucket))
			if b == nil {
				return nil
			}
			return b.ForEach(func(k, v []byte) error {
				timestamp, err := time.Parse(timeFormat, string(k))
				if err != nil {
					return err
				}
				result = append(result, entry{
					Timestamp: timestamp,
					Content:   v,
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

func (s *boltStore) GetRangeWithAggregation(parentBucketName string, start, end time.Time, agg aggregationFunction) (any, error) {
	if agg == nil {
		return nil, errors.New("aggregation function is nil")
	}
	entries, err := s.GetRange(parentBucketName, start, end)
	if err != nil {
		return nil, err
	}
	return agg(entries)
}

func (s *boltStore) Close() error {
	return s.db.Close()
}

func getBucketFromEntry(e entry) (string, error) {
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

	truncate := func(t time.Time) time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	}

	strippedStart := truncate(start)
	strippedEnd := truncate(end.AddDate(0, 0, 1))
	buckets := make([]string, 0)

	for i := strippedStart; i.Before(strippedEnd); {
		bucket, err := getBucketFromTime(i)
		if err != nil {
			return nil, err
		}
		i = i.AddDate(0, 0, 1)
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}
