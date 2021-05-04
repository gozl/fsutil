package fsutil

import (
	"os"
	"testing"
)

func TestWriteFile(t *testing.T) {
	testFile := "./testdata/write_test.txt"
	testData := "hello world"

	err := WriteFile(testFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Unexpected error writing file: %v", err)
	}

	readBytes, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Unexpected error calling os.ReadFile: %v", err)
	}

	if string(readBytes) != testData {
		t.Fatalf("Read data does not match write:\n%s\n--------\n%s", 
			testData, string(readBytes))
	}

	// write again
	testData2 := "foo bar!"
	err = WriteFile(testFile, []byte(testData2), 0644)
	if err != nil {
		t.Fatalf("Unexpected error writing file at 2nd op: %v", err)		
	}

	readBytes, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Unexpected error calling os.ReadFile: %v", err)
	}

	if string(readBytes) != testData2 {
		t.Fatalf("Read data does not match write at 2nd op:\n%s\n--------\n%s", 
			testData2, string(readBytes))
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("Unexpected error removing test file %s: %v", testFile, err)
	}
}

func TestAppendFile(t *testing.T) {
	testFile := "./testdata/append_test.txt"
	testData := "hello world"

	err := AppendFile(testFile, []byte(testData), 0644)
	if err != nil {
		t.Fatalf("Unexpected error appending file: %v", err)
	}

	readBytes, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Unexpected error calling os.ReadFile: %v", err)
	}

	if string(readBytes) != testData {
		t.Fatalf("Read data does not match write:\n%s\n--------\n%s", 
			testData, string(readBytes))
	}

	// write again
	testData2 := ", foo bar!"
	err = AppendFile(testFile, []byte(testData2), 0644)
	if err != nil {
		t.Fatalf("Unexpected error appending file at 2nd op: %v", err)		
	}

	readBytes, err = os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Unexpected error calling os.ReadFile: %v", err)
	}

	if string(readBytes) != testData+testData2 {
		t.Fatalf("Read data does not match write at 2nd op:\n%s\n--------\n%s", 
			testData2, string(readBytes))
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("Unexpected error removing test file %s: %v", testFile, err)
	}
}