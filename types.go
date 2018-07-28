package godid

import (
	"io"
	"time"
)

//entry represents one entry in the db
type entry struct {
	Timestamp time.Time
	Content   []byte
}

//aggregationFunction is a function used to aggregate entries retrieved from the store
type aggregationFunction func([]entry) (interface{}, error)

//entryStore is the db manager for entries
type entryStore interface {
	io.Closer
	Put(entry) error
	GetRange(start, end time.Time) ([]entry, error)
	GetRangeWithAggregation(start, end time.Time, agg aggregationFunction) (interface{}, error)
}
