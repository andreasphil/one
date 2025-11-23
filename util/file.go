package util

import (
	"os"
	"path/filepath"
)

func WriteTextFile(content string, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	return os.WriteFile(absPath, []byte(content), 0644)
}
