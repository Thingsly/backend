package utils

import "time"

func GetUTCTime() time.Time {
	return time.Now().UTC()
}

func GetSecondTimestamp() int64 {
	return time.Now().Unix()
}

func IsToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() &&
		t.Month() == now.Month() &&
		t.Day() == now.Day()
}

func DaysAgo(n int) time.Time {
	now := time.Now()
	past := now.AddDate(0, 0, -n)
	return past
}

func MillisecondsTimestampDaysAgo(n int) int64 {

	now := time.Now()

	past := now.AddDate(0, 0, -n)

	milliseconds := past.UnixNano() / int64(time.Millisecond)
	return milliseconds
}
