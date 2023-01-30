package main

var allowedExtensions = []string{
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

var exifFields = []string{
	"DateTimeOriginal",
	"CreationDate",
	"CreateDate",
	"ModifyDate",
}

var dateFormats = []string{
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006:01:02 15:04:05",
	"2006/01/02 15:04:05",
}
