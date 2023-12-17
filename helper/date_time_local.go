package helper

import "time"

func ParseLocal(d time.Time) time.Time {

	dateTime := time.Date(d.Year(), d.Month(), d.Day(), d.Hour(),
		d.Minute(), d.Second(), 0, time.Local)

	return dateTime
}
