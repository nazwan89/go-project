package utils

import "time"

const DefaultTimeFormat = "2006-01-02 15:04:05"

func CurrentTimestamp() string {
	loc, _ := time.LoadLocation("Asia/Kuala_Lumpur")
	return time.Now().In(loc).Format(DefaultTimeFormat)
}

func CurrentUTCTime() time.Time {
	return time.Now().UTC()
}