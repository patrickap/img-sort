package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/patrickap/img-sort/m/v2/internal/exif"
)

func IsFileExtension(path string, extensions []string) bool {
	for _, ext := range extensions {
		if strings.ToLower(filepath.Ext(path)) == ext {
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

func DecodeExif(path string) (exiftool.FileMetadata, error) {
	fileExif := exif.Instance().ExtractMetadata(path)[0]
	if fileExif.Err != nil {
		return exiftool.FileMetadata{}, fileExif.Err
	}

	return fileExif, nil
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

	return time.Time{}, errors.New("Could not parse date")
}

func ParseExifDate(fileExif exiftool.FileMetadata, exifFields, dateFormats []string) (time.Time, error) {
	var fileDate time.Time
	var fileDateErr error
	for _, exifField := range exifFields {
		if fileDate, fileDateErr = ParseDate(fileExif.Fields[exifField], dateFormats); fileDateErr == nil {
			return fileDate, nil
		}
	}

	return time.Time{}, errors.New("Could not parse exif creation date")
}

func MoveFile(path, newPath string) error {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return err
	}

	// If target file already exists, append a postfix
	if IsFileExisting(newPath) {
		fileExt := filepath.Ext(newPath)
		fileBase := filepath.Base(newPath)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			fileBaseIdx := fmt.Sprintf("%s-%d%s", fileName, i, fileExt)
			newPathIdx := filepath.Join(filepath.Dir(newPath), fileBaseIdx)
			if !IsFileExisting(newPathIdx) {
				newPath = newPathIdx
				break
			}
		}
	}

	return os.Rename(path, newPath)
}
