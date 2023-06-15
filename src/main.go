package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/barasher/go-exiftool"
)

var verison = "v0.0.5"

var versionFlag bool
var sourceFlag string
var targetFlag string
var modtimeFlag bool

var exiftoolInstance *exiftool.Exiftool

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
	var exiftoolErr error
	exiftoolInstance, exiftoolErr = exiftool.NewExiftool()
	if exiftoolErr != nil {
		fmt.Println(exiftoolErr)
		os.Exit(1)
	}
	defer exiftoolInstance.Close()

	// Recursively read source directory
	processErr := filepath.Walk(sourceFlag, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories, process files only
		if fileInfo.IsDir() {
			return nil
		}

		fmt.Printf("INF: Processing %s\n", path)

		// Allow only specified file extensions
		if !isFileExtension(path, FILE_EXTENSIONS_ALLOWED) {
			fmt.Printf("WRN: Extension %s not supported\n", filepath.Ext(path))
			return nil
		}

		// Decode file exif data and parse create date
		var fileDate time.Time
		var fileError error
		fileExif, fileError := decodeExif(path)
		fileDate, fileError = parseExifDate(fileExif, EXIF_FIELDS_DATE_CREATED, EXIF_FIELDS_DATE_FORMAT)
		if fileError != nil {
			if !modtimeFlag {
				// Move file to unknown
				newPath := filepath.Join(targetFlag, "unknown", filepath.Base(path))
				fmt.Println("WRN: Could not parse date (no modtime fallback)")
				fmt.Printf("INF: Moving to %s\n", newPath)
				return moveFile(path, newPath)
			}

			// Set file modtime as fallback
			fileDate = fileInfo.ModTime()
		}

		// Move file to destination
		yearDir := fmt.Sprintf("%d", fileDate.Year())
		monthDir := fmt.Sprintf("%d-%02d", fileDate.Year(), fileDate.Month())
		fileName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", fileDate.Year(), fileDate.Month(), fileDate.Day(), fileDate.Hour(), fileDate.Minute(), fileDate.Second(), strings.ToLower(filepath.Ext(path)))
		newPath := filepath.Join(targetFlag, yearDir, monthDir, fileName)
		fmt.Printf("INF: Moving to %s\n", newPath)
		return moveFile(path, newPath)
	})

	if processErr != nil {
		fmt.Println(processErr)
		fmt.Println("ERR: View log output above")
		os.Exit(1)
	}

	fmt.Println("INF: Completed successfully")
	os.Exit(0)
}
