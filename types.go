package godid

import "time"

//Entry represents one entry in the db
type Entry struct {
	Timestamp time.Time
	Message   []byte
}

//AggregationFunction is a function used to aggregate entries retrieved from the store
type AggregationFunction func([]Entry) (interface{}, error)

//EntryStore is the db manager for entries
type EntryStore interface {
	Put(Entry) error
	GetRange(start, end time.Time) ([]Entry, error)
	GetRangeWithAggregation(start, end time.Time, agg AggregationFunction) (interface{}, error)
}
