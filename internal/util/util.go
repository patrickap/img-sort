package util

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func IsFileExtension(extensions []string, path string) bool {
	for _, extension := range extensions {
		if strings.ToLower(filepath.Ext(path)) == extension {
			return true
		}
	}

	return false
}

func IsFileExisting(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !fileInfo.IsDir()
}

func ParseDate(dateFormats []string, dateString interface{}) (time.Time, error) {
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

func MoveFile(path string, newPath string, duplicateFileStrategy func(path string) string) error {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return err
	}

	if IsFileExisting(newPath) {
		newPath = duplicateFileStrategy(newPath)
	}

	return os.Rename(path, newPath)
}
