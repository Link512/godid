package godid

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

const (
	flatEntriesPlaceholder = "all entries"
)

var (
	lastDurationPattern = regexp.MustCompile(`^(\d+)d$`)
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

//Init initialises godid
func Init() {
	cfg, err := getConfig()
	if err != nil {
		panic(err)
	}
	if cfg == nil {
		panic(errors.New("null config"))
	}
	store, err = newBoltStore(*cfg)
	if err != nil {
		panic(err)
	}
}

// Close closes godid
func Close() {
	store.Close()
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
	start := time.Now()
	result, err := getRange(start, start, true)
	if err != nil {
		return nil, err
	}
	return result[flatEntriesPlaceholder], nil
}

//GetYesterday retrieves all entries logged yesterday
func GetYesterday() ([]string, error) {
	start := time.Now().AddDate(0, 0, -1)
	result, err := getRange(start, start, true)
	if err != nil {
		return nil, err
	}
	return result[flatEntriesPlaceholder], nil
}

//GetThisWeek returns all entries from the current week
func GetThisWeek(flat bool) (map[string][]string, error) {
	start, end := getWeekInterval(time.Now())
	return getRange(start, end, flat)
}

//GetLastWeek returns all entries from the previous week
func GetLastWeek(flat bool) (map[string][]string, error) {
	aWeekBefore := time.Now().AddDate(0, 0, -7)
	start, end := getWeekInterval(aWeekBefore)
	return getRange(start, end, flat)
}

//GetLastDuration retrives all the entries from the custom previous duration
func GetLastDuration(durationString string, flat bool) (map[string][]string, error) {

	d, err := parseDuration(durationString)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	return getRange(now.Add(-1*d), now, flat)
}

func parseDuration(durationString string) (time.Duration, error) {
	match := lastDurationPattern.FindStringSubmatch(durationString)
	if match == nil {
		return 0, didErrorf("invalid duration string %s", durationString)
	}
	d, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, err
	}
	return 24 * time.Duration(d) * time.Hour, nil
}

func getWeekInterval(reference time.Time) (time.Time, time.Time) {
	weekDay := int(reference.Weekday())
	if reference.Weekday() == time.Sunday {
		weekDay = int(time.Saturday) + 1
	}
	start := reference.AddDate(0, 0, -1*(weekDay-int(time.Monday)))
	end := reference.AddDate(0, 0, int(time.Saturday)+1-weekDay)
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
		result[flatEntriesPlaceholder] = flatEntries
	} else {
		var ok bool
		result, ok = entries.(map[string][]string)
		if !ok {
			return nil, errors.New("internal error, cannot convert result")
		}
	}
	return result, nil
}
