package util

import (
	"os"
	"path/filepath"
)

// WriteTextFile writes content to path with the given permissions. It writes to
// a temporary file in the same directory first and atomically renames it into
// place, to avoid leaving a partially written file behind in case of an error.
func WriteTextFile(content string, path string, permissions os.FileMode) error {
	success := false

	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	dir := filepath.Dir(absPath)

	tmp, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return err
	}

	defer func() {
		tmp.Close()

		if !success {
			os.Remove(tmp.Name())
		}
	}()

	if _, err := tmp.Write([]byte(content)); err != nil {
		return err
	}

	if err := tmp.Sync(); err != nil {
		return err
	}

	if err := tmp.Close(); err != nil {
		return err
	}

	if err := os.Chmod(tmp.Name(), permissions); err != nil {
		return err
	}

	if err := os.Rename(tmp.Name(), absPath); err != nil {
		return err
	}

	success = true
	return nil
}
