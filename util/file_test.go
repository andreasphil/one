package util_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/andreasphil/one/util"
)

func TestWriteTextFileWritesContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	content := "hello, world"

	if err := util.WriteTextFile(content, path, 0644); err != nil {
		t.Fatalf("WriteTextFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(got) != content {
		t.Errorf("file content = %q, want %q", got, content)
	}
}

func TestWriteTextFileSetsPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	perms := os.FileMode(0600)

	if err := util.WriteTextFile("content", path, perms); err != nil {
		t.Fatalf("WriteTextFile() error = %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	if got := info.Mode().Perm(); got != perms {
		t.Errorf("file mode = %v, want %v", got, perms)
	}
}

func TestWriteTextFileOverwritesExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")

	if err := os.WriteFile(path, []byte("old content"), 0600); err != nil {
		t.Fatalf("WriteFile() error seeding the existing file = %v", err)
	}

	newContent := "new content"
	newPerms := os.FileMode(0644)

	if err := util.WriteTextFile(newContent, path, newPerms); err != nil {
		t.Fatalf("WriteTextFile() error = %v", err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(got) != newContent {
		t.Errorf("file content = %q, want %q", got, newContent)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}

	if got := info.Mode().Perm(); got != newPerms {
		t.Errorf("file mode = %v, want %v", got, newPerms)
	}
}

func TestWriteTextFileNoLeftoverTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := util.WriteTextFile("content", path, 0644); err != nil {
		t.Fatalf("WriteTextFile() error = %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	if len(entries) != 1 || entries[0].Name() != "file.txt" {
		t.Errorf("directory contents = %v, want only file.txt", entries)
	}
}

func TestWriteTextFileErrorOnMissingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nope", "file.txt")

	if err := util.WriteTextFile("content", path, 0644); err == nil {
		t.Errorf("WriteTextFile() error = nil, want an error for the missing directory")
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("Stat() error = %v, want a not-exist error", err)
	}
}
