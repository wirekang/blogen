package par

import (
	"testing"
)

func TestDateToTime(t *testing.T) {
	s := "2020-8-12"
	ti, e := DateToTime(s)
	if e != nil {
		t.FailNow()
	}
	if ti.Day() != 12 {
		t.FailNow()
	}
	if ti.Month() != 8 {
		t.FailNow()
	}
	if ti.Year() != 2020 {
		t.FailNow()
	}
}
