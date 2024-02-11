package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadFiles(path string, include []string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(path, func(currentPath string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if len(include) == 0 || IsFileExtension(currentPath, include) {
			files = append(files, currentPath)
		}

		return nil
	})

	return files, err
}

func IsFileExtension(path string, extensions []string) bool {
	for _, extension := range extensions {
		if strings.ToLower(filepath.Ext(path)) == extension {
			return true
		}
	}

	return false
}

func IsFileExisting(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}

	return !fileInfo.IsDir()
}

func MoveFile(from string, to string) error {
	err := os.MkdirAll(filepath.Dir(to), os.ModePerm)
	if err != nil {
		return err
	}

	if IsFileExisting(to) {
		fileExt := filepath.Ext(to)
		fileBase := filepath.Base(to)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			newPath := filepath.Join(filepath.Dir(to), fmt.Sprintf("%s-%d%s", fileName, i, fileExt))
			if !IsFileExisting(newPath) {
				to = newPath
				break
			}
		}
	}

	return os.Rename(from, to)
}
