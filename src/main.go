package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/barasher/go-exiftool"
)

type FileDateTime struct {
	Type  string
	Value time.Time
}

var verison = "v0.0.3"

var versionFlag bool
var sourceFlag string
var targetFlag string
var modtimeFlag bool

var exifInstance *exiftool.Exiftool

func main() {
	flag.BoolVar(&versionFlag, "version", false, "version info")
	flag.StringVar(&sourceFlag, "source", "", "source path")
	flag.StringVar(&targetFlag, "target", "", "target path")
	flag.BoolVar(&modtimeFlag, "modtime", false, "modification time fallback")
	flag.Parse()

	if versionFlag {
		fmt.Printf("Version: %s\n", verison)
		os.Exit(0)
	}

	if sourceFlag == "" || targetFlag == "" {
		fmt.Println("Error: --source and --target are required flags.")
		fmt.Println("Usage: img-sort --source /path/to/source --target /path/to/target")
		os.Exit(1)
	}

	// Create exiftool instance
	var exifErr error
	exifInstance, exifErr = exiftool.NewExiftool()
	if exifErr != nil {
		fmt.Println(exifErr)
		os.Exit(1)
	}
	defer exifInstance.Close()

	// Recursively read source directory
	processErr := filepath.Walk(sourceFlag, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, files only
		if fileInfo.IsDir() {
			return nil
		}

		// Allow only these file extensions
		if !isExtension(path, allowedExtensions) {
			return nil
		}

		fileTime, err := decodeExifDate(path, exifFields, dateFormats)
		if err != nil {
			if modtimeFlag {
				fileTime = FileDateTime{Type: "File[ModTime]", Value: fileInfo.ModTime()}
			} else {
				return moveFileToUnknown(path, targetFlag)
			}
		}

		return moveFileToTarget(path, targetFlag, fileTime)
	})
	if processErr != nil {
		fmt.Println(processErr)
		os.Exit(1)
	}

	os.Exit(0)
}

func moveFileToUnknown(path, targetRoot string) error {
	newPath := filepath.Join(targetRoot, "unknown", filepath.Base(path))

	_, err := moveFile(path, newPath)
	if err != nil {
		return err
	}

	return nil
}

func moveFileToTarget(path string, targetRoot string, fileTime FileDateTime) error {
	yearDir := fmt.Sprintf("%d", fileTime.Value.Year())
	monthDir := fmt.Sprintf("%d-%02d", fileTime.Value.Year(), fileTime.Value.Month())
	fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileTime.Value.Year(), fileTime.Value.Month(), fileTime.Value.Day(), fileTime.Value.Hour(), fileTime.Value.Minute(), fileTime.Value.Second(), filepath.Ext(path))
	newPath := filepath.Join(targetRoot, yearDir, monthDir, fileName)

	_, err := moveFile(path, newPath)
	if err != nil {
		return err
	}

	return nil
}
