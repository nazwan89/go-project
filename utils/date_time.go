package utils

import "time"

const DefaultTimeFormat = "2006-01-02 15:04:05"

// CurrentTimestamp returns the current time formatted using the provided location.
// Pass AppConfig.Location to use the env-driven timezone.
func CurrentTimestamp(loc *time.Location) string {
	if loc == nil {
		loc = time.UTC
	}
	return time.Now().In(loc).Format(DefaultTimeFormat)
}

// CurrentUTCTime returns the current time in UTC.
// Phase 2 (OBS-04) will normalize all timestamps to RFC3339 UTC.
func CurrentUTCTime() time.Time {
	return time.Now().UTC()
}
