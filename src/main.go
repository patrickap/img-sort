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

type FileDateTime struct {
	Type  string
	Value time.Time
}

var v = "v0.0.2"
var exifInstance *exiftool.Exiftool

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

	// Create single exiftool instance
	var exifErr error
	exifInstance, exifErr = exiftool.NewExiftool()
	if exifErr != nil {
		fmt.Println(exifErr)
		os.Exit(1)
	}
	defer exifInstance.Close()

	// Recursively read source directory
	processErr := filepath.Walk(*source, func(path string, info os.FileInfo, err error) error {
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

		fmt.Printf("\nProcessing file: %s\n", filepath.Base(path))

		fileTime, err := decodeExifDateTime(path, "2006:01:02 15:04:05", *modtime)
		if err != nil {
			fmt.Println("Extracted time: <nil>")

			// If no exif data is available move file to 'unknown' directory
			newPath := filepath.Join(*target, "unknown", filepath.Base(path))

			finalPath, err := moveFile(path, newPath)
			if err != nil {
				return err
			}

			fmt.Printf("Moved file: %s\n", finalPath)
			return nil
		}

		fmt.Printf("Extracted time: %s\n", fileTime.Type)

		// Rename file and move it to target directory
		yearDir := fmt.Sprintf("%d", fileTime.Value.Year())
		monthDir := fmt.Sprintf("%d-%02d", fileTime.Value.Year(), fileTime.Value.Month())
		fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileTime.Value.Year(), fileTime.Value.Month(), fileTime.Value.Day(), fileTime.Value.Hour(), fileTime.Value.Minute(), fileTime.Value.Second(), filepath.Ext(path))
		newPath := filepath.Join(*target, yearDir, monthDir, fileName)

		finalPath, err := moveFile(path, newPath)
		if err != nil {
			return err
		}

		fmt.Printf("Moved file: %s\n", finalPath)
		return nil
	})
	if processErr != nil {
		fmt.Println(processErr)
		os.Exit(1)
	}

	fmt.Println("\nProcess completed")
	os.Exit(0)
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
	fileExif := exifInstance.ExtractMetadata(path)[0]
	if fileExif.Err != nil {
		return exiftool.FileMetadata{}, fileExif.Err
	}

	return fileExif, nil
}

func decodeExifDateTime(path, layout string, modtime bool) (FileDateTime, error) {
	fileExif, err := decodeExif(path)
	if err != nil {
		return FileDateTime{}, err
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		return FileDateTime{}, err
	}

	if fileExif.Fields["DateTimeOriginal"] != nil {
		time, err := time.Parse(layout, fileExif.Fields["DateTimeOriginal"].(string))
		if err == nil {
			return FileDateTime{Type: "EXIF[DateTimeOriginal]", Value: time}, nil
		}
	}

	if fileExif.Fields["CreateDate"] != nil {
		time, err := time.Parse(layout, fileExif.Fields["CreateDate"].(string))
		if err == nil {
			return FileDateTime{Type: "EXIF[CreateDate]", Value: time}, nil
		}
	}

	if modtime == true {
		if fileExif.Fields["ModifyDate"] != nil {
			time, err := time.Parse(layout, fileExif.Fields["ModifyDate"].(string))
			if err == nil {
				return FileDateTime{Type: "EXIF[ModifyDate]", Value: time}, nil
			}
		}

		time := fileInfo.ModTime()
		return FileDateTime{Type: "FILE[ModTime]", Value: time}, nil
	} else {
		return FileDateTime{}, errors.New("Could not decode time")
	}
}

func moveFile(path, newPath string) (string, error) {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return path, err
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
	return newPath, os.Rename(path, newPath)
}
