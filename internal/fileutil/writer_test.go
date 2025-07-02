package fileutil_test

import (
	"carte_au_tresor/internal/fileutil"
	"os"
	"testing"
)

func TestWriteLinesToFile_BasicCase(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_write_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	lines := []string{
		"C - 3 - 4",
		"M - 1 - 0",
		"T - 1 - 3 - 2",
		"A - Lara - 0 - 3 - S - 3",
	}

	// When
	err = fileutil.WriteLinesToFile(tempFile.Name(), lines)

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the file content
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := "C - 3 - 4\nM - 1 - 0\nT - 1 - 3 - 2\nA - Lara - 0 - 3 - S - 3\n"
	if string(content) != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, string(content))
	}
}

func TestWriteLinesToFile_EmptyLines(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_write_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	lines := []string{}

	// When
	err = fileutil.WriteLinesToFile(tempFile.Name(), lines)

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the file is empty
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if len(content) != 0 {
		t.Errorf("Expected empty file, got: %s", string(content))
	}
}

func TestWriteToChannel_BasicCase(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_write_channel_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	lines := make(chan string, 3)
	lines <- "C - 3 - 4"
	lines <- "M - 1 - 0"
	lines <- "T - 1 - 3 - 2"
	close(lines)

	// When
	errors := fileutil.WriteToChannel(tempFile.Name(), lines)

	// Then
	err = <-errors
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify the file content
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := "C - 3 - 4\nM - 1 - 0\nT - 1 - 3 - 2\n"
	if string(content) != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, string(content))
	}
}
