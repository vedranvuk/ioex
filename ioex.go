package ioex

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
)

// Touch updates access and modification times of specified file.
// It creates any required directories along the optionally specified path.
// If the file does not exist it is created.
// If an error occurs it is returned.
func Touch(filename string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(filename, os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

// CopyAll copies a source file or a directory to destination and does it
// recursively. Directories along the destination path(s) are created as needed.
// Files are copied using io.Copy.
//
// If overwrite is specified it silently overwrites existing destination files,
// otherwise returns an os.Exists.
//
// Permissions of source files and directories carry over to destinations.
//
// If any other error occurs is returned and it will be of type *os.PathError.
func CopyAll(destination, source string, overwrite bool) error {

	srcfi, err := os.Stat(source)
	if err != nil {
		return err
	}

	src, err := os.OpenFile(source, os.O_RDONLY, srcfi.Mode().Perm())
	if err != nil {
		return err
	}

	// Source is file. Copy and return.
	if !srcfi.IsDir() {
		flags := os.O_WRONLY | os.O_CREATE
		if !overwrite {
			flags = flags | os.O_EXCL
		}
		dst, err := os.OpenFile(destination, flags, srcfi.Mode().Perm())
		if err != nil {
			src.Close()
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			src.Close()
			dst.Close()
			return err
		}
		src.Close()
		dst.Close()
		return nil
	}

	// Create destination directory.
	if _, err := os.Stat(destination); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err := os.Mkdir(destination, srcfi.Mode().Perm()); err != nil {
			src.Close()
			return err
		}
	}

	// Enumerate files.
	infos, err := src.Readdir(-1)
	if err != nil {
		src.Close()
		return err
	}
	src.Close()

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name() < infos[j].Name()
	})

	// Recurse.
	for _, info := range infos {
		if err := CopyAll(
			filepath.Join(destination, info.Name()),
			filepath.Join(source, info.Name()),
			overwrite,
		); err != nil {
			return err
		}
	}

	return nil
}
