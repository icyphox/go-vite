package util

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyFile copies a file from src to dst.
func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// Copy modes.
	f, err := os.Stat(src)
	if err == nil {
		err = os.Chmod(dst, f.Mode())
		if err != nil {
			return err
		}
	}

	return out.Close()
}

// CopyDir copies an entire directory tree from
// src to dst.
func CopyDir(src, dst string) error {
	fi, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !fi.IsDir() {
		return fmt.Errorf("error: %q is not a directory", fi)
	}

	if err = os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	items, _ := os.ReadDir(src)
	for _, item := range items {
		srcFilename := filepath.Join(src, item.Name())
		dstFilename := filepath.Join(dst, item.Name())
		if item.IsDir() {
			if err := CopyDir(srcFilename, dstFilename); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcFilename, dstFilename); err != nil {
				return err
			}
		}
	}

	return nil
}
