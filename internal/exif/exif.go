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

	defer exif.Close()
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

func ParseDate(fileExif exiftool.FileMetadata, exifFields, dateFormats []string) (time.Time, error) {
	var fileDate time.Time
	var fileDateErr error
	for _, exifField := range exifFields {
		if fileDate, fileDateErr = util.ParseDate(fileExif.Fields[exifField], dateFormats); fileDateErr == nil {
			return fileDate, nil
		}
	}

	return time.Time{}, errors.New("failed to parse exif date")
}
