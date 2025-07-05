package common

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron"
)

func GetToday() time.Time {
	now := time.Now()
	year, month, day := now.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.Local)
}

func GetYearStart() time.Time {
	now := time.Now()
	year, _, _ := now.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
}

func GetMonthStart() time.Time {
	now := time.Now()
	year, month, _ := now.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
}

func GetYesterdayBegin() time.Time {
	return GetToday().Add(-1)
}

func DateTimeToString(date time.Time, Layout string) string {
	if Layout == "" {
		Layout = "2006-01-02 15:04:05"
	}
	return date.Format(Layout)
}

func GetWeekDay(date time.Time) int {
	currentWeekday := date.Weekday()
	var weekday int

	switch currentWeekday {
	case time.Sunday:
		weekday = 7
	case time.Monday:
		weekday = 1
	case time.Tuesday:
		weekday = 2
	case time.Wednesday:
		weekday = 3
	case time.Thursday:
		weekday = 4
	case time.Friday:
		weekday = 5
	case time.Saturday:
		weekday = 6
	}
	return weekday
}

func GetSceneExecuteTime(taskType, condition string) (time.Time, error) {
	var (
		result time.Time
		now    = time.Now()
		err    error
	)

	switch taskType {
	case "HOUR":
		// condition: "15:04:05-07:00"
		minstr := condition[:2]
		var min int
		min, err = strconv.Atoi(minstr)
		if err != nil || min > 59 || min < 0 {
			return result, errors.New("invalid time format")
		}
		if min > now.Minute() {
			result = time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), min, 0, 0, now.Location())
		} else {
			result = time.Date(now.Year(), now.Month(), now.Day(), now.Hour()+1, min, 0, 0, now.Location())
		}

	case "DAY":
		// condition: "15:04:05-07:00" 
		daytime, err := time.Parse("15:04:05-07:00", condition)
		if err != nil {
			return result, errors.New("invalid time format")
		}
		result = time.Date(now.Year(), now.Month(), now.Day(), daytime.Hour(), daytime.Minute(), daytime.Second(), 0, now.Location())
		if result.Before(now) {
			result = result.Add(time.Hour * 24)
		}

	case "WEEK":
		// condition: "1|15:04:05-07:00" 1: 1st day of the week
		parts := strings.Split(condition, "|")
		if len(parts) != 2 {
			return result, errors.New("invalid time format")
		}
		// Parse weekdays
		weekdaysStr := parts[0]
		weekdays := make([]time.Weekday, 0)
		for _, char := range weekdaysStr {
			if char >= '1' && char <= '7' {
				day, _ := strconv.Atoi(string(char))
				if day == 7 {
					day = 0
				}
				weekdays = append(weekdays, time.Weekday(day))
			}
		}
		// Parse time
		timeStr := parts[1]
		targetTime, err := time.Parse("15:04:05-07:00", timeStr)
		if err != nil {
			return result, errors.New("invalid time format")
		}
		result = GetNextTime(now, weekdays, targetTime)
	case "MONTH":
		// Parse time string
		// condition: "2T15:04:05-07:00" 2T: 2nd day of the month
		targetTime, err := time.Parse("2T15:04:05-07:00", condition)
		if err != nil {
			return result, errors.New("time parsing error")
		}
		result = getMonthNextTime(now, targetTime)
	case "CRON":
		specParser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.DowOptional | cron.Descriptor)
		schedule, err := specParser.Parse(condition)
		if err != nil {
			return result, errors.New("invalid cron format")
		}
		result = schedule.Next(now)
	default:
		return result, errors.New("unsupported time format")
	}

	return result, err
}

func GetNextTime(now time.Time, weekdays []time.Weekday, targetTime time.Time) time.Time {
	//// Get the current year, month, day, and weekday
	//year, month, day := now.Date()
	//
	//// Start searching from the current time
	//for i := 0; i < 7; i++ {
	//	// Calculate the next date that meets the condition
	//	nextDay := now.AddDate(0, 0, i)
	//	nextWeekday := nextDay.Weekday()
	//	for _, wd := range weekdays {
	//		if (wd - 1) == nextWeekday {
	//			// Set the next valid time
	//			nextTime := time.Date(year, month, day+1, targetTime.Hour(), targetTime.Minute(), targetTime.Second(), 0, time.Local)
	//			// If the time is after the current time, return this time
	//			if nextTime.After(now) {
	//				return nextTime
	//			}
	//		}
	//	}
	//}
	//
	//// If no valid time is found, return the first valid time in the next week
	//nextDay := now.AddDate(0, 0, 7-GetWeekDay(now))
	//nextWeekday := nextDay.Weekday()
	//for _, wd := range weekdays {
	//	if (wd - 1) == nextWeekday {
	//		// Set the next valid time for the found date
	//		nextTime := time.Date(year, month, day+7-GetWeekDay(now), targetTime.Hour(), targetTime.Minute(), targetTime.Second(), 0, time.Local)
	//		return nextTime
	//	}
	//}
	//return time.Time{}
	
	// Extract detailed information from the current time
	//year, month, day := now.Date()
	//hour, minute, second := now.Clock()

	// Iterate over the next 7 days to find the next valid execution time
	for i := 0; i < 7; i++ {
		// Calculate the next potential date
		nextDay := now.AddDate(0, 0, i)
		nextWeekday := nextDay.Weekday()
		for _, wd := range weekdays {
			if wd == nextWeekday {
				// Construct the target time for the current date
				nextTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), targetTime.Hour(), targetTime.Minute(), targetTime.Second(), 0, nextDay.Location())

				// If the target time is after the current time, return it
				if nextTime.After(now) {
					return nextTime
				}
			}
		}
	}

	// If no valid time is found within the next 7 days, search for the first valid time in the next week
	for i := 0; i < 7; i++ {
		// Calculate the potential date for the next week
		nextDay := now.AddDate(0, 0, 7+i)
		nextWeekday := nextDay.Weekday()

		for _, wd := range weekdays {
			if wd == nextWeekday {
				// Construct the target time for the found date in the next week
				nextTime := time.Date(nextDay.Year(), nextDay.Month(), nextDay.Day(), targetTime.Hour(), targetTime.Minute(), targetTime.Second(), 0, nextDay.Location())
				return nextTime
			}
		}
	}

	// Return zero value if no valid execution time is found
	return time.Time{}
}

// Get the next valid time that meets the conditions
func getMonthNextTime(now time.Time, targetTime time.Time) time.Time {
	// Get the current year, month, day, and hour, minute, second
	year, month, _ := now.Date()
	//targetMonth := targetTime.Month()
	targetDay := targetTime.Day()
	// Calculate the next valid time
	var nextTime time.Time
	if now.Day() <= targetDay {
		// If the current day is less than or equal to the target day, or the current month is less than the target month,
		// the next time will be the target day of the current month
		nextTime = time.Date(year, month, targetDay, targetTime.Hour(), targetTime.Minute(), targetTime.Second(), 0, time.Local)
	} else {
		// Otherwise, the next time will be the target day of the next month
		nextTime = time.Date(year, month+1, targetDay, targetTime.Hour(), targetTime.Minute(), targetTime.Second(), 0, time.Local)
	}

	// If the next time is before the current time, add one month
	if nextTime.Before(now) {
		nextTime = nextTime.AddDate(0, 1, 0)
	}

	return nextTime
}

