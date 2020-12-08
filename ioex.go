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
	var err error
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	var file *os.File
	if file, err = os.OpenFile(filename, os.O_CREATE, 0644); err != nil {
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
	var err error
	// Get source info.
	var srcinfo os.FileInfo
	if srcinfo, err = os.Stat(source); err != nil {
		return err
	}
	// Open source.
	var srcfile *os.File
	if srcfile, err = os.OpenFile(source, os.O_RDONLY, srcinfo.Mode().Perm()); err != nil {
		return err
	}
	// Source is file. Copy and return.
	if !srcinfo.IsDir() {
		var flags = os.O_WRONLY | os.O_CREATE
		if !overwrite {
			flags = flags | os.O_EXCL
		}
		var dstfile *os.File
		if dstfile, err = os.OpenFile(destination, flags, srcinfo.Mode().Perm()); err != nil {
			srcfile.Close()
			return err
		}
		if _, err = io.Copy(dstfile, srcfile); err != nil {
			srcfile.Close()
			dstfile.Close()
			return err
		}
		srcfile.Close()
		dstfile.Close()
		return nil
	}
	// Create destination directory.
	if _, err = os.Stat(destination); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		if err = os.Mkdir(destination, srcinfo.Mode().Perm()); err != nil {
			srcfile.Close()
			return err
		}
	}
	// Enumerate files.
	var infos []os.FileInfo
	if infos, err = srcfile.Readdir(-1); err != nil {
		srcfile.Close()
		return err
	}
	srcfile.Close()
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].Name() < infos[j].Name()
	})
	// Recurse.
	var info os.FileInfo
	for _, info = range infos {
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
