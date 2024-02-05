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

func Extract(paths ...string) ([]exiftool.FileMetadata, error) {
	exifs := exif.ExtractMetadata(paths...)
	isExifsErr := false
	for _, exif := range exifs {
		if exif.Err != nil {
			isExifsErr = true
		}
	}

	if isExifsErr {
		return exifs, errors.New("failed to extract exif")
	}

	return exifs, nil
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
