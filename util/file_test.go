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
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	if string(result) != content {
		t.Errorf("unexpected content: got %q, want %q", result, content)
	}
}

func TestWriteTextFileSetsPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")
	perms := os.FileMode(0600)

	if err := util.WriteTextFile("content", path, perms); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat written file: %v", err)
	}

	if info.Mode().Perm() != perms {
		t.Errorf("unexpected permissions: got %v, want %v", info.Mode().Perm(), perms)
	}
}

func TestWriteTextFileOverwritesExisting(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.txt")

	if err := os.WriteFile(path, []byte("old content"), 0600); err != nil {
		t.Fatalf("failed to seed existing file: %v", err)
	}

	newContent := "new content"
	newPerms := os.FileMode(0644)

	if err := util.WriteTextFile(newContent, path, newPerms); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}

	if string(result) != newContent {
		t.Errorf("unexpected content: got %q, want %q", result, newContent)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("failed to stat written file: %v", err)
	}

	if info.Mode().Perm() != newPerms {
		t.Errorf("unexpected permissions: got %v, want %v", info.Mode().Perm(), newPerms)
	}
}

func TestWriteTextFileNoLeftoverTempFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "file.txt")

	if err := util.WriteTextFile("content", path, 0644); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read directory: %v", err)
	}

	if len(entries) != 1 || entries[0].Name() != "file.txt" {
		t.Errorf("expected only the target file in directory, got: %v", entries)
	}
}

func TestWriteTextFileErrorOnMissingDir(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nope", "file.txt")

	if err := util.WriteTextFile("content", path, 0644); err == nil {
		t.Errorf("expected error for missing directory, got nil")
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("expected no file to be created, stat returned: %v", err)
	}
}
