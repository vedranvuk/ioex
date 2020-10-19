package ioex

import (
	"errors"
	"os"
	"testing"
)

func init() {
	os.MkdirAll("test/a/b/c", 0755)
	Touch("test/a/file.ext")
	Touch("test/a/b/file.ext")
	Touch("test/a/b/c/file.ext")

	os.Symlink("test/a/b/c", "test/link")
}

func TestCopyAll(t *testing.T) {

	// Target exists.
	if err := CopyAll("test/link", "test/a/b/c", true); !errors.Is(err, os.ErrExist) {
		t.Fatal(err)
	}

	// Successfull copy.
	if err := CopyAll("test/out", "test/a", true); !errors.Is(err, nil) {
		t.Fatal(err)
	}

	// No overwrite.
	if err := CopyAll("test/out", "test/a", false); !errors.Is(err, os.ErrExist) {
		t.Fatal(err)
	}

	// Successfull overwrite.
	if err := CopyAll("test/out", "test/a", true); !errors.Is(err, nil) {
		t.Fatal(err)
	}

}

func TestZCleanup(t *testing.T) {
	if err := os.RemoveAll("test"); err != nil {
		t.Fatal("failed to remove test data")
	}
}
