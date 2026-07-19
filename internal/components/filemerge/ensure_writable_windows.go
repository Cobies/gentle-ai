//go:build windows

package filemerge

import (
	"io/fs"
)

func ensureWritable(dir string, info fs.FileInfo, path string) error {
	// Windows directory write permissions are managed via ACLs.
	// Bypasses Unix chmod 555 checks.
	return nil
}
