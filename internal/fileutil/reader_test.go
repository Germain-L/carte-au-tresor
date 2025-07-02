package fileutil_test

import (
	"carte_au_tresor/internal/fileutil"
	"os"
	"testing"
)

func TestReadToChannel_BasicCase(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	tempFile.WriteString("C - 3 - 4\n")
	tempFile.WriteString("M - 1 - 0\n")
	tempFile.Close()

	// When
	lines, errors := fileutil.ReadToChannel(tempFile.Name())

	var result []string
	for line := range lines {
		result = append(result, line)
	}

	// Then
	select {
	case err := <-errors:
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	default:
	}
	if len(result) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(result))
	}

	if result[0] != "C - 3 - 4" {
		t.Errorf("Expected 'C - 3 - 4', got '%s'", result[0])
	}

	if result[1] != "M - 1 - 0" {
		t.Errorf("Expected 'M - 1 - 0', got '%s'", result[1])
	}
}

func TestReadToChannel_EmptyFile(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// When
	lines, errors := fileutil.ReadToChannel(tempFile.Name())

	var result []string
	for line := range lines {
		result = append(result, line)
	}

	// Then
	select {
	case err := <-errors:
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	default:
	}
	if len(result) != 0 {
		t.Errorf("Expected 0 lines, got %d", len(result))
	}
}

func TestReadToChannel_SingleLine(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	tempFile.WriteString("A - 1 - 2 - S\n")
	tempFile.Close()

	// When
	lines, errors := fileutil.ReadToChannel(tempFile.Name())

	var result []string
	for line := range lines {
		result = append(result, line)
	}

	// Then
	select {
	case err := <-errors:
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	default:
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result))
	}
	if result[0] != "A - 1 - 2 - S" {
		t.Errorf("Expected 'A - 1 - 2 - S', got '%s'", result[0])
	}
}

func TestReadToChannel_FileWithEmptyLines(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	tempFile.WriteString("C - 3 - 4\n")
	tempFile.WriteString("\n")
	tempFile.WriteString("M - 1 - 0\n")
	tempFile.WriteString("\n")
	tempFile.Close()

	// When
	lines, errors := fileutil.ReadToChannel(tempFile.Name())

	var result []string
	for line := range lines {
		result = append(result, line)
	}

	// Then
	select {
	case err := <-errors:
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	default:
	}
	if len(result) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(result))
	}
}

func TestReadToChannel_FileWithoutNewlineAtEnd(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	tempFile.WriteString("C - 3 - 4")
	tempFile.Close()

	// When
	lines, errors := fileutil.ReadToChannel(tempFile.Name())

	var result []string
	for line := range lines {
		result = append(result, line)
	}

	// Then
	select {
	case err := <-errors:
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
	default:
	}
	if len(result) != 1 {
		t.Errorf("Expected 1 line, got %d", len(result))
	}
	if result[0] != "C - 3 - 4" {
		t.Errorf("Expected 'C - 3 - 4', got '%s'", result[0])
	}
}
