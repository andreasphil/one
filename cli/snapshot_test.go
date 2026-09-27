package cli_test

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/andreasphil/one/cli"
)

const (
	sortedNotes   = "# 02.01.2026\n\nsecond\n\n# 01.01.2026\n\nfirst\n"
	unsortedNotes = "# 01.01.2026\n\nfirst\n\n# 02.01.2026\n\nsecond\n"
)

func git(t *testing.T, args ...string) string {
	t.Helper()

	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v error = %v: %s", strings.Join(args, " "), err, out)
	}

	return strings.TrimSpace(string(out))
}

// newRepo creates an empty git repository isolated from the user's git config
// and makes it the working directory for the rest of the test.
func newRepo(t *testing.T) {
	t.Helper()

	t.Setenv("GIT_CONFIG_GLOBAL", os.DevNull)
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")
	t.Setenv("GIT_AUTHOR_NAME", "Test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")

	t.Chdir(t.TempDir())
	git(t, "init", "--quiet")
}

func commitAll(t *testing.T) {
	t.Helper()

	git(t, "add", "--all")
	git(t, "commit", "--quiet", "--message", "initial")
}

func commitCount(t *testing.T) int {
	t.Helper()

	out, err := exec.Command("git", "rev-list", "--count", "HEAD").Output()
	if err != nil {
		return 0
	}

	count, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		t.Fatalf("Atoi() error = %v", err)
	}

	return count
}

func runSnapshot(t *testing.T, flags ...string) (string, error) {
	t.Helper()

	var stderr bytes.Buffer
	err := cli.Run(append([]string{"snapshot"}, flags...), io.Discard, &stderr)

	return stderr.String(), err
}

func requireFormatter(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("oxfmt"); err != nil {
		t.Skip("oxfmt is not installed")
	}
}

func TestSnapshotFailsOutsideRepository(t *testing.T) {
	t.Chdir(t.TempDir())
	writeNotes(t, "one.md", sortedNotes)

	if _, err := runSnapshot(t); err == nil {
		t.Error("snapshot() error = nil, want error")
	}
}

func TestSnapshotFailsForMissingInput(t *testing.T) {
	newRepo(t)

	if _, err := runSnapshot(t); err == nil {
		t.Error("snapshot() error = nil, want error")
	}
}

func TestSnapshotFailsWithoutChanges(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	commitAll(t)

	stderr, err := runSnapshot(t)
	if err == nil {
		t.Error("snapshot() error = nil, want error")
	}

	if !strings.Contains(stderr, "no changes since last snapshot") {
		t.Errorf("snapshot() stderr = %q, want warning about missing changes", stderr)
	}

	if got := commitCount(t); got != 1 {
		t.Errorf("commit count = %v, want 1", got)
	}
}

func TestSnapshotCommitsUntrackedInput(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)

	if stderr, err := runSnapshot(t); err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if got := git(t, "show", "HEAD:one.md") + "\n"; got != sortedNotes {
		t.Errorf("committed content = %q, want %q", got, sortedNotes)
	}
}

func TestSnapshotCommitsChangedInput(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	commitAll(t)

	changed := "# 03.01.2026\n\nthird\n\n" + sortedNotes
	writeNotes(t, "one.md", changed)

	if stderr, err := runSnapshot(t); err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if got := commitCount(t); got != 2 {
		t.Errorf("commit count = %v, want 2", got)
	}

	message := git(t, "log", "-1", "--format=%s")
	if !regexp.MustCompile(`^snapshot: \d{4}-\d{2}-\d{2} \d{2}:\d{2}$`).MatchString(message) {
		t.Errorf("commit message = %q, want snapshot: yyyy-mm-dd hh:mm", message)
	}

	if got := git(t, "show", "HEAD:one.md") + "\n"; got != changed {
		t.Errorf("committed content = %q, want %q", got, changed)
	}

	if got := git(t, "status", "--porcelain"); got != "" {
		t.Errorf("status after snapshot = %q, want clean", got)
	}
}

func TestSnapshotCommitsOnlyInput(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	writeNotes(t, "other.md", "other\n")
	commitAll(t)

	writeNotes(t, "one.md", "# 03.01.2026\n\nthird\n\n"+sortedNotes)
	writeNotes(t, "other.md", "changed\n")

	if stderr, err := runSnapshot(t); err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if got := git(t, "show", "--name-only", "--format=", "HEAD"); got != "one.md" {
		t.Errorf("committed files = %q, want one.md", got)
	}

	if got := git(t, "status", "--porcelain"); got != "M other.md" {
		t.Errorf("status after snapshot = %q, want other.md modified", got)
	}
}

func TestSnapshotAllowsStagedInput(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	commitAll(t)

	writeNotes(t, "one.md", "# 03.01.2026\n\nthird\n\n"+sortedNotes)
	git(t, "add", "one.md")

	if stderr, err := runSnapshot(t); err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if got := commitCount(t); got != 2 {
		t.Errorf("commit count = %v, want 2", got)
	}
}

func TestSnapshotFailsWhenOtherFilesAreStaged(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	writeNotes(t, "other.md", "other\n")
	commitAll(t)

	writeNotes(t, "one.md", "# 03.01.2026\n\nthird\n\n"+sortedNotes)
	writeNotes(t, "other.md", "changed\n")
	git(t, "add", "other.md")

	_, err := runSnapshot(t)
	if err == nil || !strings.Contains(err.Error(), "other files are currently staged for commit") {
		t.Errorf("snapshot() error = %v, want error about staged files", err)
	}

	if got := commitCount(t); got != 1 {
		t.Errorf("commit count = %v, want 1", got)
	}

	if got := git(t, "diff", "--cached", "--name-only"); got != "other.md" {
		t.Errorf("staged files = %q, want other.md", got)
	}
}

func TestSnapshotDetectsStagedFilesOutsideWorkingDirectory(t *testing.T) {
	newRepo(t)
	if err := os.Mkdir("notes", 0755); err != nil {
		t.Fatalf("Mkdir() error = %v", err)
	}
	writeNotes(t, "notes/one.md", sortedNotes)
	writeNotes(t, "other.md", "other\n")
	commitAll(t)

	writeNotes(t, "notes/one.md", "# 03.01.2026\n\nthird\n\n"+sortedNotes)
	writeNotes(t, "other.md", "changed\n")
	git(t, "add", "other.md")
	t.Chdir("notes")

	if _, err := runSnapshot(t); err == nil {
		t.Error("snapshot() error = nil, want error about staged files")
	}
}

func TestSnapshotAcceptsAbsoluteInput(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	commitAll(t)

	writeNotes(t, "one.md", "# 03.01.2026\n\nthird\n\n"+sortedNotes)
	git(t, "add", "one.md")

	input, err := filepath.Abs("one.md")
	if err != nil {
		t.Fatalf("Abs() error = %v", err)
	}

	if stderr, err := runSnapshot(t, "--input", input); err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if got := git(t, "show", "--name-only", "--format=", "HEAD"); got != "one.md" {
		t.Errorf("committed files = %q, want one.md", got)
	}
}

func TestSnapshotDetectsStagedFilesWithAbsoluteInput(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	writeNotes(t, "other.md", "other\n")
	commitAll(t)

	writeNotes(t, "one.md", "# 03.01.2026\n\nthird\n\n"+sortedNotes)
	writeNotes(t, "other.md", "changed\n")
	git(t, "add", "other.md")

	input, err := filepath.Abs("one.md")
	if err != nil {
		t.Fatalf("Abs() error = %v", err)
	}

	_, err = runSnapshot(t, "--input", input)
	if err == nil || !strings.Contains(err.Error(), "other files are currently staged for commit") {
		t.Errorf("snapshot() error = %v, want error about staged files", err)
	}
}

func TestSnapshotCheckMakesNoChanges(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	commitAll(t)

	writeNotes(t, "one.md", "# 03.01.2026\n\nthird\n\n"+sortedNotes)

	if stderr, err := runSnapshot(t, "--check"); err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if got := commitCount(t); got != 1 {
		t.Errorf("commit count = %v, want 1", got)
	}

	if got := git(t, "status", "--porcelain"); got != "M one.md" {
		t.Errorf("status after check = %q, want one.md modified and unstaged", got)
	}
}

func TestSnapshotCheckFailsWithoutChanges(t *testing.T) {
	newRepo(t)
	writeNotes(t, "one.md", sortedNotes)
	commitAll(t)

	if _, err := runSnapshot(t, "--check"); err == nil {
		t.Error("snapshot() error = nil, want error")
	}
}

func TestSnapshotTidySortsBeforeCommitting(t *testing.T) {
	requireFormatter(t)
	newRepo(t)
	writeNotes(t, "one.md", unsortedNotes)

	stderr, err := runSnapshot(t, "--tidy")
	if err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if !strings.Contains(stderr, "sorted") || !strings.Contains(stderr, "already formatted") {
		t.Errorf("snapshot() stderr = %q, want sorted and already formatted", stderr)
	}

	if got := git(t, "show", "HEAD:one.md") + "\n"; got != sortedNotes {
		t.Errorf("committed content = %q, want %q", got, sortedNotes)
	}

	if got := git(t, "status", "--porcelain"); got != "" {
		t.Errorf("status after snapshot = %q, want clean", got)
	}
}

func TestSnapshotTidyFormatsBeforeCommitting(t *testing.T) {
	requireFormatter(t)
	newRepo(t)
	writeNotes(t, "one.md", "# 02.01.2026\n\n* second\n\n# 01.01.2026\n\nfirst\n")

	stderr, err := runSnapshot(t, "--tidy")
	if err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if !strings.Contains(stderr, "formatted") || !strings.Contains(stderr, "already sorted") {
		t.Errorf("snapshot() stderr = %q, want formatted and already sorted", stderr)
	}

	if got := git(t, "show", "HEAD:one.md") + "\n"; got != "# 02.01.2026\n\n- second\n\n# 01.01.2026\n\nfirst\n" {
		t.Errorf("committed content = %q, want formatted notes", got)
	}
}

func TestSnapshotCheckWithTidyMakesNoChanges(t *testing.T) {
	requireFormatter(t)
	newRepo(t)
	writeNotes(t, "one.md", unsortedNotes)

	stderr, err := runSnapshot(t, "--check", "--tidy")
	if err != nil {
		t.Fatalf("snapshot() error = %v, stderr = %v", err, stderr)
	}

	if !strings.Contains(stderr, "notes would be sorted") {
		t.Errorf("snapshot() stderr = %q, want notes would be sorted", stderr)
	}

	content, err := os.ReadFile("one.md")
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if string(content) != unsortedNotes {
		t.Errorf("content after check = %q, want unchanged", content)
	}

	if got := commitCount(t); got != 0 {
		t.Errorf("commit count = %v, want 0", got)
	}
}
