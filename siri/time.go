package siri

import "time"

// Time helpers

// ParseISOTime attempts to parse common ISO-8601/RFC3339 timestamps with offsets.
func ParseISOTime(s string) (time.Time, bool) {
	if s == "" {
		return time.Time{}, false
	}
	layouts := []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05-07:00",
		"2006-01-02T15:04:05Z07:00",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

// Latest returns the latest of t1 and t2, handling zero values.
func Latest(t1, t2 time.Time) time.Time {
	if t1.IsZero() {
		return t2
	}
	if t2.IsZero() {
		return t1
	}
	if t2.After(t1) {
		return t2
	}
	return t1
}

// FormatDateYYYYMMDD formats a time into YYYYMMDD string.
func FormatDateYYYYMMDD(t time.Time) string { return t.Format("20060102") }
