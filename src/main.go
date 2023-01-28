package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
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

	err := filepath.Walk(*source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !isValidFileType(path) {
			return nil
		}

		fileTime, err := extractExifDateTime(path)
		if err != nil {
			newPath := filepath.Join(*target, "unknown", filepath.Base(path))
			return moveFile(path, newPath)
		}

		year := fmt.Sprintf("%d", fileTime.Year())
		month := fmt.Sprintf("%d-%02d", fileTime.Year(), fileTime.Month())
		fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileTime.Year(), fileTime.Month(), fileTime.Day(), fileTime.Hour(), fileTime.Minute(), fileTime.Second(), filepath.Ext(path))
		newPath := filepath.Join(*target, year, month, fileName)

		return moveFile(path, newPath)

	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isValidFileType(path string) bool {
	ext := filepath.Ext(path)
	validTypes := []string{".tiff", ".tif", ".gif", ".jpeg", ".jpg", ".png", ".img", ".bmp", ".raw", ".heif", ".heic", ".mkv", ".avi", ".mov", ".wmv", ".mp4", ".m4v", ".mpg", ".mpeg", ".hevc"}

	for _, t := range validTypes {
		if ext == t {
			return true
		}
	}

	return false
}

func extractExifDateTime(path string) (time.Time, error) {
	f, err := os.Open(path)
	if err != nil {
		return time.Time{}, err
	}
	defer f.Close()

	x, err := exif.Decode(f)
	if err != nil {
		return time.Time{}, err
	}

	tm, err := x.DateTime()
	if err != nil {
		return time.Time{}, err
	}

	return tm, nil
}

func moveFile(sourceFile, targetFile string) error {
	err := os.MkdirAll(filepath.Dir(targetFile), os.ModePerm)
	if err != nil {
		return err
	}

	// Check if target file already exists
	if _, err := os.Stat(targetFile); err == nil {
		// Target file already exists, append a postfix
		ext := filepath.Ext(targetFile)
		base := targetFile[:len(targetFile)-len(ext)]
		for i := 1; ; i++ {
			newTargetFile := fmt.Sprintf("%s-%d%s", base, i, ext)
			if _, err := os.Stat(newTargetFile); os.IsNotExist(err) {
				targetFile = newTargetFile
				break
			}
		}
	}

	return os.Rename(sourceFile, targetFile)
}

// func isFile(path string) bool {
// 	info, err := os.Stat(path)
// 	if err != nil {
// 		return false
// 	}

// 	return !info.IsDir()
// }
