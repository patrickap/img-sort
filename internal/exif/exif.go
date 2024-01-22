package exif

import (
	"github.com/barasher/go-exiftool"
	"github.com/patrickap/img-sort/m/v2/internal/log"
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
