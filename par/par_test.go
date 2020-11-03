package par

import (
	"fmt"
	"testing"
)

func TestExtract(t *testing.T) {
	str := []byte("key1: val1\nkey2: val2\n$\n# this is markdown")
	set, md, err := Extract(str, '#')
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("==%s==\n==%s==\n", set, md)
}

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
