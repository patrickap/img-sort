package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

func isExtAllowed(path string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.ToLower(filepath.Ext(path)) == ext {
			return true
		}
	}

	return false
}

func isFileExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !fileInfo.IsDir()
}

func parseDate(dateString interface{}, dateFormats []string) (time.Time, error) {
	isString := dateString != nil && reflect.TypeOf(dateString).Kind() == reflect.String
	if isString {
		for _, format := range dateFormats {
			date, err := time.Parse(format, dateString.(string))
			if err == nil {
				return date, nil
			}
		}
	}

	return time.Time{}, errors.New("Could not parse date")
}

func decodeExif(path string) (exiftool.FileMetadata, error) {
	fileExif := exiftoolInstance.ExtractMetadata(path)[0]
	if fileExif.Err != nil {
		return exiftool.FileMetadata{}, fileExif.Err
	}

	return fileExif, nil
}

func moveFile(path, newPath string) error {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return err
	}

	// If target file already exists, append a postfix
	if isFileExist(newPath) {
		fileExt := filepath.Ext(newPath)
		fileBase := filepath.Base(newPath)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			fileBaseIdx := fmt.Sprintf("%s-%d%s", fileName, i, fileExt)
			newPathIdx := filepath.Join(filepath.Dir(newPath), fileBaseIdx)
			if !isFileExist(newPathIdx) {
				newPath = newPathIdx
				break
			}
		}
	}

	// Transform file base and extension to lowercase
	newPath = filepath.Join(filepath.Dir(newPath), strings.ToLower(filepath.Base(newPath)))
	return os.Rename(path, newPath)
}
