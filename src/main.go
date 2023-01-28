package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	source := flag.String("source", "", "source path")
	target := flag.String("target", "", "target path")

	flag.Parse()

	if *source == "" || *target == "" {
		fmt.Println("Error: --source and --target are required flags.")
		fmt.Println("Usage: img-sort --source /path/to/source --target /path/to/target")
		os.Exit(1)
	}

	// Recursively read source directory
	err := filepath.Walk(*source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, files only
		if info.IsDir() {
			return nil
		}

		// Allow valid file extensions only
		if !isValidExt(path) {
			return nil
		}

		fileTime, err := decodeExifTime(path)
		// If no exif data is available move file to 'unknown' directory
		if err != nil {
			newPath := filepath.Join(*target, "unknown", filepath.Base(path))
			return moveFile(path, newPath)
		}

		// Rename file and move it to target directory
		yearDir := fmt.Sprintf("%d", fileTime.Year())
		monthDir := fmt.Sprintf("%d-%02d", fileTime.Year(), fileTime.Month())
		fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileTime.Year(), fileTime.Month(), fileTime.Day(), fileTime.Hour(), fileTime.Minute(), fileTime.Second(), filepath.Ext(path))
		newPath := filepath.Join(*target, yearDir, monthDir, fileName)

		return moveFile(path, newPath)

	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isValidExt(path string) bool {
	validExts := []string{".tiff", ".tif", ".gif", ".jpeg", ".jpg", ".png", ".img", ".bmp", ".raw", ".heif", ".heic", ".mkv", ".avi", ".mov", ".wmv", ".mp4", ".m4v", ".mpg", ".mpeg", ".hevc"}

	for _, ext := range validExts {
		if filepath.Ext(path) == ext {
			return true
		}
	}

	return false
}

func decodeExif(path string) (*exif.Exif, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return exif.Decode(file)
}

func decodeExifTime(path string) (time.Time, error) {
	exifData, err := decodeExif(path)
	if err != nil {
		return time.Time{}, err
	}

	exifTime, err := exifData.DateTime()
	if err != nil {
		return time.Time{}, err
	}

	return exifTime, nil
}

func moveFile(sourceFile, targetFile string) error {
	err := os.MkdirAll(filepath.Dir(targetFile), os.ModePerm)
	if err != nil {
		return err
	}

	// If target file already exists, append a postfix
	if _, err := os.Stat(targetFile); err == nil {
		fileExt := filepath.Ext(targetFile)
		fileBase := filepath.Base(targetFile)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			newTargetFile := fmt.Sprintf("%s-%d%s", fileName, i, fileExt)
			if _, err := os.Stat(newTargetFile); os.IsNotExist(err) {
				targetFile = newTargetFile
				break
			}
		}
	}

	return os.Rename(sourceFile, targetFile)
}
