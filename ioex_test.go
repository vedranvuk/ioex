package ioex

import (
	"errors"
	"os"
	"testing"
)

func createTestData() {
	var err error
	if err = os.MkdirAll("test/a/b/c", 0755); err != nil {
		panic(err)
	}
	Touch("test/a/file.ext")
	Touch("test/a/b/file.ext")
	Touch("test/a/b/c/file.ext")

	os.Symlink("test/a/b/c", "test/link")
}

func deleteTestData() {
	if err := os.RemoveAll("test"); err != nil {
		panic("failed to remove test data.")
	}
}

func TestExists(t *testing.T) {
	createTestData()
	var exists bool
	var err error
	if exists, err = Exists("test/a/file.ext"); err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("Exists failed.")
	}
	if exists, err = Exists("test/doesnotexist.exe"); err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatal("Exists failed.")
	}
	deleteTestData()
}

func TestTouch(t *testing.T) {
	var err error
	if err = Touch("test/sub/index.html"); err != nil {
		t.Fatal(err)
	}
	var exists bool
	if exists, err = Exists("test/sub/index.html"); err != nil {
		t.Fatal(err)
	}
	if !exists {
		t.Fatal("Touch failed.")
	}
	deleteTestData()
}

func TestCopyAll(t *testing.T) {
	createTestData()
	defer deleteTestData()
	// Exists.
	if err := CopyAll("test/link", "test/a/b/c", true); !errors.Is(err, os.ErrExist) {
		t.Fatal(err)
	}
	// Successfull copy.
	if err := CopyAll("test/out", "test", true); !errors.Is(err, nil) {
		t.Fatal(err)
	}
	// No overwrite.
	if err := CopyAll("test/out", "test/a", false); err != nil {
		t.Fatal(err)
	}
	// Successfull overwrite.
	if err := CopyAll("test/out", "test/a", true); err != nil {
		t.Fatal(err)
	}
}
