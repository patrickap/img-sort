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

var v = "v0.0.1"

func main() {
	version := flag.Bool("version", false, "version info")
	source := flag.String("source", "", "source path")
	target := flag.String("target", "", "target path")

	flag.Parse()

	if *version != false {
		fmt.Printf("Version: %s\n", v)
		os.Exit(0)
	}

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
		if !isExtensionValid(path) {
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

func isExtensionValid(path string) bool {
	extensions := []string{".tiff", ".tif", ".gif", ".jpeg", ".jpg", ".png", ".img", ".bmp", ".raw", ".heif", ".heic", ".mkv", ".avi", ".mov", ".wmv", ".mp4", ".m4v", ".mpg", ".mpeg", ".hevc"}

	for _, ext := range extensions {
		if strings.ToLower(filepath.Ext(path)) == ext {
			return true
		}
	}

	return false
}

func isFileExisting(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !info.IsDir()
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
		info, err := os.Stat(path)
		if err != nil {
			return time.Time{}, err
		}

		return info.ModTime(), err
	}

	return exifData.DateTime()
}

func moveFile(path, newPath string) error {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return err
	}

	// If target file already exists, append a postfix
	if isFileExisting(newPath) {
		fileExt := filepath.Ext(newPath)
		fileBase := filepath.Base(newPath)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			fileBaseIdx := fmt.Sprintf("%s-%d%s", fileName, i, fileExt)
			newPathIdx := filepath.Join(filepath.Dir(newPath), fileBaseIdx)
			if !isFileExisting(newPathIdx) {
				newPath = newPathIdx
				break
			}
		}
	}

	// Transform file base and extension to lowercase
	newPath = filepath.Join(filepath.Dir(newPath), strings.ToLower(filepath.Base(newPath)))

	fmt.Printf("Move: %s -> %s\n", filepath.Base(path), newPath)
	return os.Rename(path, newPath)
}
