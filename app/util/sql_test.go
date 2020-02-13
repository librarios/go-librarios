package util

import (
	"github.com/google/go-cmp/cmp"
	"gopkg.in/guregu/null.v3"
	"testing"
	"time"
)

func TestNewNullTime(t *testing.T) {
	locSeoul, _ := time.LoadLocation("Asia/Seoul")

	testCases := []struct {
		source   string
		expected null.Time
	}{
		{"", null.TimeFromPtr(nil)},
		{"foo", null.TimeFromPtr(nil)},
		{"2014-11-17T09:10:11.123+09:00", null.TimeFrom(time.Date(2014, time.November, 17, 9, 10, 11, 123*1e6, locSeoul))},
	}

	for _, testCase := range testCases {
		actual := NullTimeFromString(testCase.source)
		if diff := cmp.Diff(testCase.expected, actual); diff != "" {
			t.Errorf("TestNewNullTime() mismatch (-expected +actual):\n%s", diff)
		}
	}
}
