package util

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadFiles(path string, includeExtensions []string) ([]string, error) {
	files := []string{}
	err := filepath.WalkDir(path, func(currentPath string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil
		}

		if IsFileExtension(includeExtensions, path) {
			files = append(files, currentPath)
		}

		return nil
	})

	return files, err
}

func IsFileExtension(extensions []string, path string) bool {
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

func MoveFile(path string, newPath string) error {
	err := os.MkdirAll(filepath.Dir(newPath), os.ModePerm)
	if err != nil {
		return err
	}

	if IsFileExisting(newPath) {
		fileExt := filepath.Ext(newPath)
		fileBase := filepath.Base(newPath)
		fileName := strings.TrimSuffix(fileBase, fileExt)

		for i := 1; ; i++ {
			newPathIdx := filepath.Join(filepath.Dir(newPath), fmt.Sprintf("%s-%d%s", fileName, i, fileExt))
			if !IsFileExisting(newPathIdx) {
				newPath = newPathIdx
				break
			}
		}
	}

	return os.Rename(path, newPath)
}
