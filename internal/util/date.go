package util

import (
	"errors"
	"reflect"
	"time"
)

func TryParseDate(dateFormats []string, dateString interface{}) (time.Time, error) {
	isString := dateString != nil && reflect.TypeOf(dateString).Kind() == reflect.String
	if isString {
		for _, format := range dateFormats {
			date, err := time.Parse(format, dateString.(string))
			if err == nil {
				return date, nil
			}
		}
	}

	return time.Time{}, errors.New("failed to parse date")
}
