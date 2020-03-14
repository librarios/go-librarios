package util

import (
	"github.com/google/go-cmp/cmp"
	"testing"
	"time"
)

func TestParseTime(t *testing.T) {
	data := "2014-11-17T09:10:11.123+09:00"
	loc, _ := time.LoadLocation("Asia/Seoul")
	expected := time.Date(2014, time.November, 17, 9, 10, 11, 123*1e6, loc)
	actual, err := ParseTime(data)

	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("TestParseTime() mismatch (-expected +actual):\n%s", diff)
	}
}
