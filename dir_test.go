package fsutil

import (
	"fmt"
	"os"
	"testing"
)

func TestIsEmptyDir(t *testing.T) {
	testDir := "./testdata/foo_dir"
	testFile := "./testdata/foo_dir/bar.txt"
	testData := "hello world"

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Unexpected error creating test dir: %v", err)
	}

	ok, err := IsEmptyDir(testDir)
	if err != nil {
		t.Fatalf("Unexpected error reading empty dir: %v", err)
	}
	if !ok {
		t.Fatalf("Expect dir to be empty")
	}

	err = WriteFile(testFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Unexpected error writing file: %v", err)
	}

	ok, err = IsEmptyDir(testDir)
	if err != nil {
		t.Fatalf("Unexpected error reading empty dir on 2nd op: %v", err)
	}
	if ok {
		t.Fatalf("Expect dir to be not empty")
	}

	ok, err = IsEmptyDir(testFile)
	if err == nil {
		t.Fatalf("Expect non-nil error reading file as dir")
	}
	if ok {
		t.Fatalf("Expect false reading file as dir")
	}
	if err != ErrNotDir {
		t.Fatalf("Expected error to be ErrNotDir but got: %v", err)
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("Unexpected error removing file: %v", err)
	}

	ok, err = IsEmptyDir(testDir)
	if err != nil {
		t.Fatalf("Unexpected error reading empty dir on 3rd op: %v", err)
	}
	if !ok {
		t.Fatalf("Expect dir to be empty")
	}

	err = os.Remove(testDir)
	if err != nil {
		t.Fatalf("Unexpected error removing empty dir: %v", err)
	}
}

func TestDir(t *testing.T) {
	testRootDir := "./testdata/a_dir"
	testDir := testRootDir + "/bar_dir"
	testFile := testRootDir + "/bar.txt"

	err := os.MkdirAll(testDir, 0755)
	if err != nil {
		t.Fatalf("Unexpected error creating test dir: %v", err)
	}

	err = WriteFile(testFile, []byte("hello world"), 0644)
	if err != nil {
		t.Fatalf("Unexpected error writing file: %v", err)
	}

	x, err := Dir(testRootDir, "", 0)
	if err != nil {
		t.Fatalf("Unexpected error dir: %v", err)
	}

	for _, v := range x {
		fmt.Println(v)
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("Unexpected error removing file: %v", err)
	}

	err = os.Remove(testDir)
	if err != nil {
		t.Fatalf("Unexpected error removing dir: %v", err)
	}

	err = os.Remove(testRootDir)
	if err != nil {
		t.Fatalf("Unexpected error removing dir: %v", err)
	}
}