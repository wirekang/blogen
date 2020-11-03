package sett

import (
	"testing"
)

func TestSett(t *testing.T) {
	str := "KEY1: value1\n" +
		"key2 : 32\n \n\n\t\n" +
		"Key_3: one, two , three\n\n \n" +
		"4key  :false\n" +
		"#Key5:true\n"
	ss, err := ParseSettings(str)
	if err != nil {
		t.Fatal(err)
	}
	find := func(key string) Setting {
		s, err := ss.Find(key)
		if err != nil {
			t.Fatal(err)
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
	if i, err := key2.IntValue(); i != 32 || err != nil {
		t.Fatal(err)
	}

	ar := key3.StringArrayValue()
	if ar[0] != "one" || ar[1] != "two" || ar[2] != "three" {
		t.FailNow()
	}
	if b, err := key4.BoolValue(); b || err != nil {
		t.Fatal(err)
	}

}
