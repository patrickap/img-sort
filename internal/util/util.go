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

func ReadFiles(path string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(path, func(currentPath string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		files = append(files, currentPath)
		return nil
	})

	return files, err
}

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

func MoveFile(path string, newPath string) error {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return err
	}

	if IsFileExisting(newPath) {
		fileExt := filepath.Ext(newPath)
		fileBase := filepath.Base(newPath)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			newPathIdx := filepath.Join(filepath.Dir(newPath), fmt.Sprintf("%s-%d%s", fileName, i, fileExt))
			if !IsFileExisting(newPathIdx) {
				newPath = newPathIdx
				break
			}
		}
	}

	return os.Rename(path, newPath)
}
