package exif

import (
	"errors"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/patrickap/img-sort/m/v2/internal/log"
	"github.com/patrickap/img-sort/m/v2/internal/util"
)

var (
	exif    *exiftool.Exiftool
	exifErr error
)

func init() {
	exif, exifErr = exiftool.NewExiftool()
	if exifErr != nil {
		log.Error().Msg("Failed to initialize exiftool")
		panic(exifErr)
	}
}

func Instance() *exiftool.Exiftool {
	return exif
}

func Decode(path string) (exiftool.FileMetadata, error) {
	fileExif := exif.ExtractMetadata(path)[0]
	if fileExif.Err != nil {
		return exiftool.FileMetadata{}, fileExif.Err
	}

	return fileExif, nil
}

func ParseDate(dateFormats []string, exifFields []string, fileExif exiftool.FileMetadata) (time.Time, error) {
	for _, exifField := range exifFields {
		fileDate, err := util.ParseDate(dateFormats, fileExif.Fields[exifField])
		if err == nil {
			return fileDate, nil
		}
	}

	return time.Time{}, errors.New("failed to parse exif date")
}
