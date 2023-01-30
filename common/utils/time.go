package utils

import (
	"github.com/dustin/go-humanize"
	"math"
	"time"
)

var (
	customTimeMagnitudes = []humanize.RelTimeMagnitude{
		{D: time.Second, Format: "now", DivBy: time.Second},
		{D: time.Minute, Format: "%ds %s", DivBy: time.Second},
		{D: time.Hour, Format: "%dm %s", DivBy: time.Minute},
		{D: humanize.Day, Format: "%dh %s", DivBy: time.Hour},
		{D: humanize.Week, Format: "%dd %s", DivBy: humanize.Day},
		{D: humanize.Year, Format: "%dw %s", DivBy: humanize.Week},
		{D: humanize.LongTime, Format: "%dy %s", DivBy: humanize.Year},
		{D: math.MaxInt64, Format: "a long while %s", DivBy: 1},
	}
	lastseenTimeMagnitudes = []humanize.RelTimeMagnitude{
		{D: humanize.Day, Format: "today", DivBy: time.Hour},
		{D: humanize.Week, Format: "%dd %s", DivBy: humanize.Day},
		{D: humanize.Year, Format: "%dw %s", DivBy: humanize.Week},
		{D: humanize.LongTime, Format: "%dy %s", DivBy: humanize.Year},
		{D: math.MaxInt64, Format: "never", DivBy: 1},
	}
)

func CustomRelTime(a, b time.Time, albl, blbl string) string {
	return humanize.CustomRelTime(a, b, albl, blbl, customTimeMagnitudes)
}
func CustomTime(then time.Time) string {
	return CustomRelTime(then, time.Now(), "ago", "from now")
}
func LastSeenRelTime(a, b time.Time, albl, blbl string) string {
	return humanize.CustomRelTime(a, b, albl, blbl, lastseenTimeMagnitudes)
}
func LastSeenTime(then time.Time) string {
	return LastSeenRelTime(then, time.Now(), "ago", "from now")
}

func GetFirstDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 0, -t.Day()+1)
}
func GetLastDayOfMonth(t time.Time) time.Time {
	return t.AddDate(0, 1, -t.Day())
}

// GetFormatted20060102 Returns time as "2006-01-02" format, hence other remaining values are zero
func GetFormatted20060102(t time.Time) (time.Time, error) {
	newTime, err := time.Parse("2006-01-02", t.Format("2006-01-02"))
	if err != nil {
		return t, err
	}
	return newTime, nil
}
