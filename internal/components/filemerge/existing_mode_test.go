package filemerge

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestExistingFileMode(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not meaningful on Windows")
	}
	dir := t.TempDir()
	private := filepath.Join(dir, "private.json")
	if err := os.WriteFile(private, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(private, 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(dir, "link.json")
	if err := os.Symlink(private, link); err != nil {
		t.Fatal(err)
	}

	for _, tt := range []struct {
		name string
		path string
		want fs.FileMode
	}{
		{name: "existing file keeps its mode", path: private, want: 0o600},
		{name: "missing file uses fallback", path: filepath.Join(dir, "missing.json"), want: 0o644},
		{name: "directory uses fallback", path: dir, want: 0o644},
		{name: "symlink uses fallback", path: link, want: 0o644},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExistingFileMode(tt.path, 0o644); got != tt.want {
				t.Fatalf("ExistingFileMode(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestExistingFileModePreservedByWriteFileAtomic(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are not meaningful on Windows")
	}
	path := filepath.Join(t.TempDir(), "opencode.json")
	if err := os.WriteFile(path, []byte("{}"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := WriteFileAtomic(path, []byte(`{"share":"disabled"}`), ExistingFileMode(path, 0o644)); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode after rewrite = %v, want 0600", got)
	}
}
