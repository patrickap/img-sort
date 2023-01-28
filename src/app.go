package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

func main() {
	sourcePath, targetPath := parseCommandLine()

	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if !validFileType(path) {
			return nil
		}

		fileTime, err := extractExifDateTime(path)
		if err != nil {
			fmt.Println("FILE NOK", path)
			return moveFileToUnknown(path, targetPath)
		}

		fmt.Println("FILE OK", path)
		return moveFileToTarget(path, targetPath, fileTime)
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func parseCommandLine() (string, string) {
	args := os.Args[1:]

	if len(args) != 4 || args[0] != "--source" || args[2] != "--target" {
		fmt.Println("Usage: sort --source /path/to/source --target /path/to/target")
		os.Exit(1)
	}

	return args[1], args[3]
}

func validFileType(path string) bool {
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

// TODO: single method for moving to target or unknown, also if duplicate add postfix

func moveFileToTarget(path, targetPath string, tm time.Time) error {
	year := fmt.Sprintf("%d", tm.Year())
	yearPath := filepath.Join(targetPath, year)

	month := fmt.Sprintf("%d-%02d", tm.Year(), tm.Month())
	monthPath := filepath.Join(yearPath, month)

	newName := fmt.Sprintf("%d-%02d-%02d_%02d.%02d.%02d%s", tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), filepath.Ext(path))
	newPath := filepath.Join(monthPath, newName)
	newDir := filepath.Dir(newPath)

	err := os.MkdirAll(newDir, fs.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("NEW PATH", newPath)
	err = os.Rename(path, newPath)
	if err != nil {
		return err
	}

	return nil
}

func moveFileToUnknown(path, targetPath string) error {
	name := filepath.Base(path)
	newPath := filepath.Join(targetPath, "unknown", name)
	newDir := filepath.Dir(newPath)

	err := os.MkdirAll(newDir, fs.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("NEW PATH", newDir)
	err = os.Rename(path, newPath)
	if err != nil {
		return err
	}

	return nil
}
