//Package sett parses the following format:
//
// ------------------------------
//
// key: string value
//
// key: 32
//
// key: true
//
// key: string, array, value
//
// #comment
//
// ------------------------------
package sett

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/wirekang/blogen/er"
)

// Setting contains key and parsable value
type Setting struct {
	key   string
	value string
}

// Settings is array of Setting.
// you can find a Setting by Find(key) function.
type Settings []Setting

// Find returns a Setting with matching key.
func (ss Settings) Find(key string) (Setting, error) {
	for _, s := range ss {
		if s.key == key {
			return s, nil
		}
	}
	return Setting{}, errors.New("no match key")
}

// ParseSettings parses Settings from string
func ParseSettings(str string) (Settings, error) {
	lines := strings.Split(str, "\n")
	settings := make([]Setting, 0, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[0] == '#' {
			continue
		}
		kv := strings.Split(line, ":")
		if len(kv) != 2 {
			return nil, fmt.Errorf("\"%s\" is not key:value format", line)
		}

		key := strings.TrimSpace(kv[0])
		if key == "" {
			return nil, fmt.Errorf("\"%s\" is not key", key)
		}
		value := strings.TrimSpace(kv[1])
		if value == "" {
			return nil, fmt.Errorf("\"%s\" is not value", value)
		}
		key = strings.ToLower(key)
		settings = append(settings, Setting{key: key, value: value})
	}
	return settings, nil
}

// ParseSettingsFromFile calls ParseSettings
func ParseSettingsFromFile(file string) ([]Setting, error) {
	bytes, err := ioutil.ReadFile(file)
	if er.PrintIfNotNil(err) {
		return nil, nil
	}
	return ParseSettings(string(bytes))
}

// Key returns key of setting
func (s Setting) Key() string {
	return s.key
}

// StringValue returns value as string
func (s Setting) StringValue() string {
	return s.value
}

// BoolValue returns value as bool
func (s Setting) BoolValue() (bool, error) {
	return strconv.ParseBool(s.value)
}

// IntValue returns value as int
func (s Setting) IntValue() (int, error) {
	return strconv.Atoi(s.value)
}

// StringArrayValue returns value as array of string
func (s Setting) StringArrayValue() []string {
	ar := strings.Split(s.value, ",")
	re := make([]string, len(ar))
	for i, s := range ar {
		re[i] = strings.TrimSpace(s)
	}
	return re
}

// MakePlaceholder returns setting format string with empty values.
//
// - - -
//
// key1:
//
// key2:
//
// key3:
func MakePlaceholder(keys []string) string {
	str := ""
	for _, key := range keys {
		str += key + ": \n"
	}
	return str
}

// MakePlaceholderFile calls MakePlaceholder and writes result to file.
func MakePlaceholderFile(file string, keys []string) error {
	return ioutil.WriteFile(file, []byte(MakePlaceholder(keys)), 0755)
}
