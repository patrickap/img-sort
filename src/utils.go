package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

func isExtension(path string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.ToLower(filepath.Ext(path)) == ext {
			return true
		}
	}

	return false
}

func isExisting(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
}

func parseDate(date string, formats []string) (time.Time, error) {
	for _, format := range formats {
		t, err := time.Parse(format, date)
		if err == nil {
			return t, nil
		}
	}

	return time.Time{}, errors.New("Could not parse date using provided formats")
}

func decodeExif(path string) (exiftool.FileMetadata, error) {
	fileExif := exifInstance.ExtractMetadata(path)[0]
	if fileExif.Err != nil {
		return exiftool.FileMetadata{}, fileExif.Err
	}

	return fileExif, nil
}

func decodeExifDate(path string, fields []string, formats []string) (FileDateTime, error) {
	fileExif, err := decodeExif(path)
	if err != nil {
		return FileDateTime{}, err
	}

	for _, field := range fields {
		fieldValue := fileExif.Fields[field]
		if fieldValue != nil {
			date, err := parseDate(fieldValue.(string), formats)
			if err != nil {
				continue
			}

			return FileDateTime{Type: fmt.Sprintf("Exif[%s]", field), Value: date}, nil
		}
	}

	return FileDateTime{}, errors.New("Could not decode date using provided fields and formats")
}

func moveFile(path, newPath string) (string, error) {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return path, err
	}

	// If target file already exists, append a postfix
	if isExisting(newPath) {
		fileExt := filepath.Ext(newPath)
		fileBase := filepath.Base(newPath)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			fileBaseIdx := fmt.Sprintf("%s-%d%s", fileName, i, fileExt)
			newPathIdx := filepath.Join(filepath.Dir(newPath), fileBaseIdx)
			if !isExisting(newPathIdx) {
				newPath = newPathIdx
				break
			}
		}
	}

	// Transform file base and extension to lowercase
	newPath = filepath.Join(filepath.Dir(newPath), strings.ToLower(filepath.Base(newPath)))
	return newPath, os.Rename(path, newPath)
}
