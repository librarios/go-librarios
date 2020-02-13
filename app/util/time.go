package util

import "time"

func ParseTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
