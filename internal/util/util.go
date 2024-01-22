package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"
)

func IsExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.ToLower(filepath.Ext(file)) == ext {
			return true
		}
	}

	return false
}

func IsExisting(file string) bool {
	fileInfo, err := os.Stat(file)
	if os.IsNotExist(err) {
		return false
	}

	return !fileInfo.IsDir()
}

func ParseDate(dateString interface{}, dateFormats []string) (time.Time, error) {
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

func Move(path, newPath string) error {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return err
	}

	// If target file already exists, append a postfix
	if IsExisting(newPath) {
		fileExt := filepath.Ext(newPath)
		fileBase := filepath.Base(newPath)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			fileBaseIdx := fmt.Sprintf("%s-%d%s", fileName, i, fileExt)
			newPathIdx := filepath.Join(filepath.Dir(newPath), fileBaseIdx)
			if !IsExisting(newPathIdx) {
				newPath = newPathIdx
				break
			}
		}
	}

	return os.Rename(path, newPath)
}
