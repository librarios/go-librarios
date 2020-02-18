package util

import (
	"gopkg.in/guregu/null.v3"
	"strings"
)

func NullTimeFromString(s string) null.Time {
	if t, err := ParseTime(s); err != nil {
		return null.TimeFromPtr(nil)
	} else {
		return null.TimeFrom(t)
	}
}

func NullString(s string) null.String {
	if s == "" {
		return null.StringFromPtr(nil)
	} else {
		return null.StringFrom(s)
	}
}

func NullStringJoin(arr []string, separator string) null.String {
	if len(arr) == 0 {
		return null.StringFromPtr(nil)
	} else {
		return null.StringFrom(strings.Join(arr, separator))
	}
}

func NullFloat(f float64) null.Float {
	return null.NewFloat(f, f != 0)
}

func NullInt(i int64) null.Int {
	return null.NewInt(i, i != 0)
}