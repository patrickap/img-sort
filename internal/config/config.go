package config

var FILE_EXTENSIONS_SUPPORTED = []string{
	".tiff",
	".tif",
	".gif",
	".jpeg",
	".jpg",
	".png",
	".img",
	".bmp",
	".raw",
	".heif",
	".heic",
	".mkv",
	".avi",
	".mov",
	".wmv",
	".mp4",
	".m4v",
	".mpg",
	".mpeg",
	".hevc",
}

// The order determines which exif field is tried to be read first
// The availability of exif information and the exact field names
// may vary depending on the specific file format and software used to create it.
var EXIF_FIELDS_DATE_CREATED = []string{
	// tiff, tif, jpeg, jpg, img, raw, heif, heic, hevc
	"DateTimeOriginal",
	// mkv, avi, mov, wmv, mp4, m4v, mpg, mpeg
	"CreationDate",
	// png, img
	"CreationTime",
	// heif, heic
	"CreateDate",
}

var EXIF_FIELDS_DATE_FORMAT = []string{
	// Default exif date format
	"2006:01:02 15:04:05",
	"2006:01:02 15:04:05Z",
	"2006:01:02 15:04:05+07:00",
	"2006:01:02 15:04:05-07:00",

	"2006-01-02T15:04:05",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05+07:00",
	"2006-01-02T15:04:05-07:00",

	"2006-01-02 15:04:05",
	"2006-01-02 15:04:05Z",
	"2006-01-02 15:04:05+07:00",
	"2006-01-02 15:04:05-07:00",

	"2006/01/02 15:04:05",
	"2006/01/02 15:04:05Z",
	"2006/01/02 15:04:05+07:00",
	"2006/01/02 15:04:05-07:00",
}
