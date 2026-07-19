//go:build !windows

package filemerge

import (
	"fmt"
	"io/fs"
	"os"
)

func ensureWritable(dir string, info fs.FileInfo, path string) error {
	if info.Mode().Perm()&0o200 == 0 {
		if err := os.Chmod(dir, 0o755); err != nil {
			return fmt.Errorf("relax parent directory permissions for %q: %w", path, err)
		}
	}
	return nil
}
