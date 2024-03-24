package godid

import (
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

const (
	flatEntriesPlaceholder = "all entries"
	rootBucketName         = "root"
)

var (
	lastDurationPattern = regexp.MustCompile(`^(\d+)d$`)
)

var (
	store           entryStore
	flatAggregation = func(entries []entry) (any, error) {
		return lo.Map(entries, func(e entry, _ int) string {
			return string(e.Content)
		}), nil
	}
	perDayAggregation = func(entries []entry) (any, error) {
		result := make(map[string][]string)
		for _, entry := range entries {
			content := string(entry.Content)
			bucket, err := getBucketFromEntry(entry)
			if err != nil {
				return nil, err
			}
			result[bucket] = append(result[bucket], content)
		}
		return result, nil
	}
)

// Init initialises godid
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

// AddEntry adds an entry to the underlying store in the root bucket
func AddEntry(what string) error {
	return AddEntryToBucket(rootBucketName, what)
}

// AddEntryToBucket adds an entry to the underlying store in the specified parent bucket
func AddEntryToBucket(bucket string, what string) error {
	e := entry{
		Content:   []byte(what),
		Timestamp: time.Now(),
	}
	err := store.Put(bucket, e)
	if err != nil {
		getLogger().WithFields(logrus.Fields{
			"component": "manager",
			"method":    "AddEntry",
			"entry":     what,
		}).WithError(err).Error("failed to put entry")
	}
	return err
}

// GetToday retrieves all entries logged today from the root bucket
func GetToday() ([]string, error) {
	return GetTodayFromBucket(rootBucketName)
}

// GetTodayFromBucket retrieves all entries logged today from the specified bucket
func GetTodayFromBucket(bucketName string) ([]string, error) {
	start := time.Now()
	result, err := getRange(bucketName, start, start, true)
	if err != nil {
		getLogger().WithFields(logrus.Fields{
			"component": "manager",
			"method":    "GetToday",
			"timestamp": start,
		}).WithError(err).Error("failed to get entries")
		return nil, err
	}
	return result[flatEntriesPlaceholder], nil
}

// GetYesterday retrieves all entries logged yesterday from the root bucket
func GetYesterday() ([]string, error) {
	return GetYesterdayFromBucket(rootBucketName)
}

// GetYesterdayFromBucket retrieves all entries logged yesterday from the specified bucket
func GetYesterdayFromBucket(bucketName string) ([]string, error) {
	start := time.Now().AddDate(0, 0, -1)
	result, err := getRange(bucketName, start, start, true)
	if err != nil {
		getLogger().WithFields(logrus.Fields{
			"component": "manager",
			"method":    "GetYesterday",
			"timestamp": start,
		}).WithError(err).Error("failed to get entries")
		return nil, err
	}
	return result[flatEntriesPlaceholder], nil
}

// GetThisWeek returns all entries from the current week from the root bucket
func GetThisWeek(flat bool) (map[string][]string, error) {
	return GetThisWeekFromBucket(rootBucketName, flat)
}

// GetThisWeekFromBucket returns all entries from the current week from the specified bucket
func GetThisWeekFromBucket(bucketName string, flat bool) (map[string][]string, error) {
	start, end := getWeekInterval(time.Now())
	result, err := getRange(bucketName, start, end, flat)
	if err != nil {
		getLogger().WithFields(logrus.Fields{
			"component": "manager",
			"method":    "GetThisWeek",
			"start":     start,
			"end":       end,
			"flat":      flat,
		}).WithError(err).Error("failed to get entries")
	}
	return result, err
}

// GetLastWeek returns all entries from the previous week from the root bucket
func GetLastWeek(flat bool) (map[string][]string, error) {
	return GetLastWeekFromBucket(rootBucketName, flat)
}

// GetLastWeekFromBucket returns all entries from the previous week from the specified bucket
func GetLastWeekFromBucket(bucketName string, flat bool) (map[string][]string, error) {
	aWeekBefore := time.Now().AddDate(0, 0, -7)
	start, end := getWeekInterval(aWeekBefore)
	result, err := getRange(bucketName, start, end, flat)
	if err != nil {
		getLogger().WithFields(logrus.Fields{
			"component": "manager",
			"method":    "GetThisWeek",
			"start":     start,
			"end":       end,
			"flat":      flat,
		}).WithError(err).Error("failed to get entries")
	}
	return result, err
}

// GetLastDuration retrives all the entries from the custom previous duration from the root bucket
func GetLastDuration(durationString string, flat bool) (map[string][]string, error) {
	return GetLastDurationFromBucket(rootBucketName, durationString, flat)
}

// GetLastDurationFromBucket retrives all the entries from the custom previous duration from the specified bucket
func GetLastDurationFromBucket(bucketName string, durationString string, flat bool) (map[string][]string, error) {

	d, err := parseDuration(durationString)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	result, err := getRange(bucketName, now.Add(-1*d), now, flat)
	if err != nil {
		getLogger().WithFields(logrus.Fields{
			"component": "manager",
			"method":    "GetThisWeek",
			"start":     now.Add(-1 * d),
			"end":       now,
			"flat":      flat,
		}).WithError(err).Error("failed to get entries")
	}
	return result, err
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

func getRange(bucketName string, start, end time.Time, flat bool) (map[string][]string, error) {
	var agg aggregationFunction
	if flat {
		agg = flatAggregation
	} else {
		agg = perDayAggregation
	}

	entries, err := store.GetRangeWithAggregation(bucketName, start, end, agg)
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
