package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

var v = "v0.0.2"

func main() {
	version := flag.Bool("version", false, "version info")
	source := flag.String("source", "", "source path")
	target := flag.String("target", "", "target path")
	modtime := flag.Bool("modtime", false, "modification time fallback")

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

	err := process(*source, *target, *modtime)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func process(source, target string, modtime bool) error {
	// Recursively read source directory
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
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

		fmt.Printf("\nProcessing file %s\n", path)

		fileTime, err := decodeExifDateTime(path, "2006:01:02 15:04:05", modtime)
		// If no exif data is available move file to 'unknown' directory
		if err != nil {
			newPath := filepath.Join(target, "unknown", filepath.Base(path))

			return moveFile(path, newPath)
		}

		// Rename file and move it to target directory
		yearDir := fmt.Sprintf("%d", fileTime.Year())
		monthDir := fmt.Sprintf("%d-%02d", fileTime.Year(), fileTime.Month())
		fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileTime.Year(), fileTime.Month(), fileTime.Day(), fileTime.Hour(), fileTime.Minute(), fileTime.Second(), filepath.Ext(path))
		newPath := filepath.Join(target, yearDir, monthDir, fileName)

		return moveFile(path, newPath)
	})

	return err
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

func decodeExif(path string) (exiftool.FileMetadata, error) {
	exif, err := exiftool.NewExiftool()
	if err != nil {
		return exiftool.FileMetadata{}, err
	}
	defer exif.Close()

	return exif.ExtractMetadata(path)[0], nil
}

func decodeExifDateTime(path, layout string, modtime bool) (time.Time, error) {
	fmt.Printf("Decoding exif datetime of %s\n", filepath.Base(path))

	fileExif, err := decodeExif(path)
	if err != nil {
		return time.Time{}, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return time.Time{}, err
	}

	if fileExif.Fields["DateTimeOriginal"] != nil {
		date, err := time.Parse(layout, fileExif.Fields["DateTimeOriginal"].(string))
		if err == nil {
			fmt.Println("Using DateTimeOriginal (exif)")
			return date, nil
		}
	}

	if fileExif.Fields["CreateDate"] != nil {
		date, err := time.Parse(layout, fileExif.Fields["CreateDate"].(string))
		if err == nil {
			fmt.Println("Using CreateDate (exif)")
			return date, nil
		}
	}

	if modtime == true {
		if fileExif.Fields["ModifyDate"] != nil {
			date, err := time.Parse(layout, fileExif.Fields["ModifyDate"].(string))
			if err == nil {
				fmt.Println("Fallback to ModifyDate (exif)")
				return date, nil
			}
		}

		fmt.Println("Fallback to ModTime (file)")
		return fileInfo.ModTime(), nil
	} else {
		fmt.Println("No datetime available")
		return time.Time{}, errors.New("Missing datetime")
	}
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

	fmt.Printf("Move file to %s\n", newPath)
	return os.Rename(path, newPath)
}
