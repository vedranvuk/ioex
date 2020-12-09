package ioex

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Exists returns if the file specified by filename exists.
// If an error occurs it is returned and the exists result is invalid.
func Exists(filename string) (bool, error) {
	var err error
	if _, err = os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

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

// CopyAll copies from source to destination where destination must be a
// directory if it exists or will be created as directory if it does not.
// If it exists and is not a directory CopyAll returns os.ErrExist.
//
// Source can be a file or a directory. If it is a file it is copied to
// destination. If it is a directory it is enumerated and all its children are
// copied to destination creating subdirectories at destination as needed.
//
// Symbolic links are skipped silently. Copying one is a noop and returns nil.
//
// If overwrite is specified it silently overwrites existing destination files,
// otherwise returns an os.Exists mid operation with incomplete copy results.
//
// Permissions of source files and directories carry over to destinations.
//
// If any other error occurs is returned and it will be of type *os.PathError.
func CopyAll(destination, source string, overwrite bool) error {
	var err error
	// Get source info.
	var srcinfo os.FileInfo
	if srcinfo, err = os.Lstat(source); err != nil {
		return err
	}
	// Skip symlinks.
	if srcinfo.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	// Source is a file.
	if !srcinfo.IsDir() {
		return CopyFile(destination, source, overwrite)
	}
	// Source is directory, enumerate it.
	var infos []os.FileInfo
	if infos, err = ioutil.ReadDir(source); err != nil {
		return err
	}
	var info os.FileInfo
	// Recurse children.
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

// CopyFile copies source file to destination file along with permissions or
// returns an error. Symlinks are skipped and return a nil error.
func CopyFile(destination, source string, overwrite bool) error {
	var err error
	// Get source info.
	var srcinfo os.FileInfo
	if srcinfo, err = os.Lstat(source); err != nil {
		return err
	}

	// Must be a file.
	if srcinfo.IsDir() {
		return errors.New("ioex: source is directory")
	}
	// Must not be a symlink.
	if srcinfo.Mode()&os.ModeSymlink != 0 {
		return nil
	}
	// Get destination parent directory info.
	var dstinfo os.FileInfo
	if dstinfo, err = os.Lstat(filepath.Dir(destination)); err != nil {
		// If destination parent does not exist create it.
		if !os.IsNotExist(err) {
			return err
		}
		// Get source parent directory info.
		if dstinfo, err = os.Lstat(filepath.Dir(source)); err != nil {
			return err
		}
		// Create destination parent dir.
		if err = os.MkdirAll(filepath.Dir(destination), dstinfo.Mode().Perm()); err != nil {
			return err
		}
		// Get destination parent directory info.
		if dstinfo, err = os.Lstat(filepath.Dir(destination)); err != nil {
			return err
		}
	}
	// Destination parent must be a dir if it exists.
	if !dstinfo.IsDir() {
		return os.ErrExist
	}
	// Open source file.
	var srcfile *os.File
	if srcfile, err = os.OpenFile(source, os.O_RDONLY, srcinfo.Mode().Perm()); err != nil {
		return err
	}
	defer srcfile.Close()
	// Open destination file.
	var flags = os.O_WRONLY | os.O_CREATE
	if !overwrite {
		flags = flags | os.O_EXCL
	}
	var dstfile *os.File
	if dstfile, err = os.OpenFile(
		destination,
		flags,
		srcinfo.Mode().Perm(),
	); err != nil {
		return err
	}
	defer dstfile.Close()
	if _, err = io.Copy(dstfile, srcfile); err != nil {
		return err
	}
	return nil
}
