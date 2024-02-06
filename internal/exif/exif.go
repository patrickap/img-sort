package exif

import (
	"errors"
	"time"

	"github.com/barasher/go-exiftool"
	"github.com/patrickap/img-sort/m/v2/internal/config"
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

func Extract(paths ...string) []exiftool.FileMetadata {
	return exif.ExtractMetadata(paths...)
}

func ParseDate(exifFields []string, fileExif exiftool.FileMetadata) (time.Time, error) {
	for _, exifField := range exifFields {
		date, err := util.ParseDate(config.EXIF_FIELDS_DATE_FORMAT, fileExif.Fields[exifField])
		if err == nil {
			return date, nil
		}
	}

	return time.Time{}, errors.New("failed to parse exif date")
}
