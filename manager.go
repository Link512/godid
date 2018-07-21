package godid

import (
	"errors"
	"time"
)

var (
	store           entryStore
	flatAggregation = func(entries []entry) (interface{}, error) {
		result := make([]string, len(entries))
		for i, entry := range entries {
			result[i] = string(entry.Content)
		}
		return result, nil
	}
	perDayAggregation = func(entries []entry) (interface{}, error) {
		result := make(map[string][]string)
		for _, entry := range entries {
			content := string(entry.Content)
			var toAdd []string
			bucket, err := getBucketFromEntry(entry)
			if err != nil {
				return nil, err
			}
			if e, found := result[bucket]; found {
				toAdd = append(e, content)
			} else {
				toAdd = []string{content}
			}
			result[bucket] = toAdd
		}
		return result, nil
	}
)

func init() {
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	store, err = newBoltStore(*cfg)
	if err != nil {
		panic(err)
	}
}

//AddEntry adds an entry to the underlying store
func AddEntry(what string) error {
	e := entry{
		Content:   []byte(what),
		Timestamp: time.Now(),
	}
	return store.Put(e)
}

//GetToday retrieves all entries logged today
func GetToday() ([]string, error) {
	end := time.Now()
	start := end.Add(-1 * time.Hour)
	result, err := getRange(start, end, true)
	if err != nil {
		return nil, err
	}
	return result["flat"], nil
}

//GetLastWeek returns all entries from the previous week
func GetLastWeek(flat bool) (map[string][]string, error) {
	start, end := getPreviousWeekInterval(time.Now())
	return getRange(start, end, flat)
}

//GetLastDuration retrives all the entries from the custom previous duration
func GetLastDuration(durationString string, flat bool) (map[string][]string, error) {
	d, err := time.ParseDuration(durationString)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return getRange(now.Add(-1*d), now, flat)
}

func getPreviousWeekInterval(reference time.Time) (time.Time, time.Time) {
	sub := -7
	if reference.Weekday() == time.Sunday {
		sub--
	}
	aWeekBefore := reference.AddDate(0, 0, sub)
	start := aWeekBefore.AddDate(0, 0, -1*(int(aWeekBefore.Weekday())-int(time.Monday)))
	end := aWeekBefore.AddDate(0, 0, int(time.Saturday)+1-int(aWeekBefore.Weekday()))

	return start, end
}

func getRange(start, end time.Time, flat bool) (map[string][]string, error) {
	var agg aggregationFunction
	if flat {
		agg = flatAggregation
	} else {
		agg = perDayAggregation
	}

	entries, err := store.GetRangeWithAggregation(start, end, agg)
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)

	if flat {
		flatEntries, ok := entries.([]string)
		if !ok {
			return nil, errors.New("internal error, cannot convert result")
		}
		result["flat"] = flatEntries
	} else {
		var ok bool
		result, ok = entries.(map[string][]string)
		if !ok {
			return nil, errors.New("internal error, cannot convert result")
		}
	}
	return result, nil
}
