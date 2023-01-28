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
	source := flag.String("source", "", "the source path")
	target := flag.String("target", "", "the target path")

	flag.Parse()

	if *source == "" || *target == "" {
		fmt.Println("Error: --source and --target are required flags.")
		fmt.Println("Usage: go run main.go --source /path/to/source --target /path/to/target")
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
			newPath := filepath.Join(*target, "unknown")
			return moveFile(path, newPath)
		}

		year := fmt.Sprintf("%d", fileTime.Year())
		month := fmt.Sprintf("%d-%02d", fileTime.Year(), fileTime.Month())
		timestamp := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileTime.Year(), fileTime.Month(), fileTime.Day(), fileTime.Hour(), fileTime.Minute(), fileTime.Second(), filepath.Ext(path))
		newPath := filepath.Join(*target, year, month, timestamp)

		return moveFile(path, newPath)

	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func isValidFileType(path string) bool {
	ext := filepath.Ext(path)
	validTypes := []string{".jpg", ".png", ".mov", ".mp4"}

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

func createDir(path string) error {
	if isFile(path) {
		// If the path is a file path, use the parent directory
		path = filepath.Dir(path)
	}

	// Create the directory, including all its parent directories if missing
	return os.MkdirAll(path, os.ModePerm)
}

func moveFile(src, dst string) error {
	if isDir(src) {
		return fmt.Errorf("Cannot move file. Source '%s' is a directory.", src)
	}

	if isDir(dst) {
		// Destination is a directory, append the file name
		dst = filepath.Join(dst, filepath.Base(src))
	}

	err := createDir(dst)
	if err != nil {
		return err
	}

	// Check if target file already exists
	if _, err := os.Stat(dst); err == nil {
		// Target file already exists, append a postfix
		ext := filepath.Ext(dst)
		base := dst[:len(dst)-len(ext)]
		for i := 1; ; i++ {
			newDst := fmt.Sprintf("%s-%d%s", base, i, ext)
			if _, err := os.Stat(newDst); os.IsNotExist(err) {
				dst = newDst
				break
			}
		}
	}

	return os.Rename(src, dst)
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.IsDir()
}
