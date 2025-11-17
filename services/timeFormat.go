package services

import "time"

const TimeStringFormat string = "02.01.2006"

func TimeFormat(t time.Time) string {
	return t.Format(TimeStringFormat)
}
