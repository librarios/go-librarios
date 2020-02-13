package util

import (
	"gopkg.in/guregu/null.v3"
)

func NullTimeFromString(s string) null.Time {
	if t, err := ParseTime(s); err != nil {
		return null.TimeFromPtr(nil)
	} else {
		return null.TimeFrom(t)
	}
}
