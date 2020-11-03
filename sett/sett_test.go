package sett

import (
	"testing"

	"github.com/wirekang/blogen/er"
)

func TestSett(t *testing.T) {
	str := "KEY1: value1\n" +
		"key2 : 32\n \n\n\t\n" +
		"Key_3: one, two , three\n\n \n" +
		"4key  :false\n" +
		"#Key5:true\n"
	ss, err := ParseSettings(str)
	if er.PrintIfNotNil(err) {
		t.FailNow()
	}
	find := func(key string) Setting {
		s, err := ss.Find(key)
		if er.PrintIfNotNil(err) {
			t.FailNow()
		}
		return s
	}
	key1 := find("key1")
	key2 := find("key2")
	key3 := find("key_3")
	key4 := find("4key")
	_, err = ss.Find("#key5")
	if err == nil {
		t.FailNow()
	}

	if key1.StringValue() != "value1" {
		t.FailNow()
	}
	if i, err := key2.IntValue(); i != 32 || er.PrintIfNotNil(err) {
		t.FailNow()
	}

	ar := key3.StringArrayValue()
	if ar[0] != "one" || ar[1] != "two" || ar[2] != "three" {
		t.FailNow()
	}
	if b, err := key4.BoolValue(); b || er.PrintIfNotNil(err) {
		t.FailNow()
	}

}
