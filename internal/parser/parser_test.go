package parser_test

import (
	"carte_au_tresor/internal/parser"
	"os"
	"strings"
	"testing"
)

func TestParseCoords(t *testing.T) {
	parser := parser.NewParser()

	tests := []struct {
		input       string
		expectedRow int
		expectedCol int
		expectError bool
	}{
		{"C - 3 - 4", 4, 3, false},
		{"M - 1 - 0", 0, 1, false},
		{"T - 0 - 3 - 2", 3, 0, false},
		{"A - Lara - 1 - 1 - S - AADADAGGA", 1, 1, false},
		{"# This is a comment", 0, 0, true},
		{"", 0, 0, true},
		{"Invalid line", 0, 0, true},
		{"X - abc - def", 0, 0, true},
	}

	for _, test := range tests {
		// Since parseCoords is private, we'll test it through public methods that use it
		t.Run(test.input, func(t *testing.T) {
			if strings.HasPrefix(test.input, "C") {
				_, _, err := parser.ParseMapDimensions(test.input)
				if test.expectError && err == nil {
					t.Errorf("Expected error for input %q, but got none", test.input)
				}
				if !test.expectError && err != nil {
					t.Errorf("Unexpected error for input %q: %v", test.input, err)
				}
			}
		})
	}
}

func TestParseMapDimensions(t *testing.T) {
	parser := parser.NewParser()

	tests := []struct {
		input          string
		expectedWidth  int
		expectedHeight int
		expectError    bool
	}{
		{"C - 3 - 4", 3, 4, false},
		{"C - 10 - 5", 10, 5, false},
		{"C - 0 - 0", 0, 0, false},
		{"C - abc - def", 0, 0, true},
		{"Invalid", 0, 0, true},
		{"# Comment", 0, 0, true},
	}

	for _, test := range tests {
		width, height, err := parser.ParseMapDimensions(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", test.input, err)
			}
			if width != test.expectedWidth {
				t.Errorf("Expected width %d, got %d for input %q", test.expectedWidth, width, test.input)
			}
			if height != test.expectedHeight {
				t.Errorf("Expected height %d, got %d for input %q", test.expectedHeight, height, test.input)
			}
		}
	}
}

func TestParseMountain(t *testing.T) {
	parser := parser.NewParser()

	tests := []struct {
		input       string
		expectedRow int
		expectedCol int
		expectError bool
	}{
		{"M - 1 - 0", 0, 1, false},
		{"M - 2 - 1", 1, 2, false},
		{"M - 0 - 0", 0, 0, false},
		{"M - abc - def", 0, 0, true},
		{"Invalid", 0, 0, true},
		{"# Comment", 0, 0, true},
	}

	for _, test := range tests {
		mountain, err := parser.ParseMountain(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", test.input, err)
			}
			if mountain.Row != test.expectedRow {
				t.Errorf("Expected row %d, got %d for input %q", test.expectedRow, mountain.Row, test.input)
			}
			if mountain.Col != test.expectedCol {
				t.Errorf("Expected col %d, got %d for input %q", test.expectedCol, mountain.Col, test.input)
			}
		}
	}
}

func TestParseTreasure(t *testing.T) {
	parser := parser.NewParser()

	tests := []struct {
		input         string
		expectedRow   int
		expectedCol   int
		expectedCount int
		expectError   bool
	}{
		{"T - 0 - 3 - 2", 3, 0, 2, false},
		{"T - 1 - 3 - 3", 3, 1, 3, false},
		{"T - 5 - 5 - 1", 5, 5, 1, false},
		{"T - 0 - 0 - 0", 0, 0, 0, false},
		{"T - 1 - 2", 0, 0, 0, true},     // Missing treasure count
		{"T - abc - def - 2", 0, 0, 0, true}, // Invalid coordinates
		{"T - 1 - 2 - abc", 0, 0, 0, true},   // Invalid treasure count
		{"Invalid", 0, 0, 0, true},
		{"# Comment", 0, 0, 0, true},
	}

	for _, test := range tests {
		treasure, err := parser.ParseTreasure(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", test.input, err)
			}
			if treasure.Row != test.expectedRow {
				t.Errorf("Expected row %d, got %d for input %q", test.expectedRow, treasure.Row, test.input)
			}
			if treasure.Col != test.expectedCol {
				t.Errorf("Expected col %d, got %d for input %q", test.expectedCol, treasure.Col, test.input)
			}
			if treasure.Count != test.expectedCount {
				t.Errorf("Expected count %d, got %d for input %q", test.expectedCount, treasure.Count, test.input)
			}
		}
	}
}

func TestParseAdventurer(t *testing.T) {
	parser := parser.NewParser()

	tests := []struct {
		input             string
		expectedName      string
		expectedRow       int
		expectedCol       int
		expectedDirection string
		expectedActions   string
		expectError       bool
	}{
		{"A - Lara - 1 - 1 - S - AADADAGGA", "Lara", 1, 1, "S", "AADADAGGA", false},
		{"A - Indiana - 0 - 0 - N - ADAGA", "Indiana", 0, 0, "N", "ADAGA", false},
		{"A - Test - 5 - 3 - E - A", "Test", 3, 5, "E", "A", false},
		{"A - Lara", "", 0, 0, "", "", true},                     // Missing fields
		{"A - abc - def - 1 - S - ADAGA", "", 0, 0, "", "", true}, // Invalid coordinates
		{"Invalid", "", 0, 0, "", "", true},
		{"# Comment", "", 0, 0, "", "", true},
	}

	for _, test := range tests {
		adventurer, err := parser.ParseAdventurer(test.input)

		if test.expectError {
			if err == nil {
				t.Errorf("Expected error for input %q, but got none", test.input)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", test.input, err)
			}
			if adventurer.Name != test.expectedName {
				t.Errorf("Expected name %q, got %q for input %q", test.expectedName, adventurer.Name, test.input)
			}
			if adventurer.Row != test.expectedRow {
				t.Errorf("Expected row %d, got %d for input %q", test.expectedRow, adventurer.Row, test.input)
			}
			if adventurer.Col != test.expectedCol {
				t.Errorf("Expected col %d, got %d for input %q", test.expectedCol, adventurer.Col, test.input)
			}
			if adventurer.Direction != test.expectedDirection {
				t.Errorf("Expected direction %q, got %q for input %q", test.expectedDirection, adventurer.Direction, test.input)
			}
			if adventurer.Actions != test.expectedActions {
				t.Errorf("Expected actions %q, got %q for input %q", test.expectedActions, adventurer.Actions, test.input)
			}
			if len(adventurer.Treasures) != 0 {
				t.Errorf("Expected empty treasures array, got %v for input %q", adventurer.Treasures, test.input)
			}
		}
	}
}

func TestLoadGameData(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_game_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	gameContent := `C - 3 - 4
M - 1 - 0
M - 2 - 1
T - 0 - 3 - 2
T - 1 - 3 - 3
A - Lara - 1 - 1 - S - AADADAGGA`

	tempFile.WriteString(gameContent)
	tempFile.Close()

	// When
	gameParser := parser.NewParser()
	mapData, err := gameParser.LoadGameData(tempFile.Name())

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mapData.Width != 3 || mapData.Height != 4 {
		t.Errorf("Expected map dimensions 3x4, got %dx%d", mapData.Width, mapData.Height)
	}

	if len(mapData.Mountains) != 2 {
		t.Errorf("Expected 2 mountains, got %d", len(mapData.Mountains))
	}

	if len(mapData.Treasures) != 2 {
		t.Errorf("Expected 2 treasures, got %d", len(mapData.Treasures))
	}

	if len(mapData.Adventurers) != 1 {
		t.Errorf("Expected 1 adventurer, got %d", len(mapData.Adventurers))
	}

	// Check adventurer details
	adventurer := mapData.Adventurers[0]
	if adventurer.Name != "Lara" {
		t.Errorf("Expected adventurer name 'Lara', got '%s'", adventurer.Name)
	}

	if adventurer.Direction != "S" {
		t.Errorf("Expected adventurer direction 'S', got '%s'", adventurer.Direction)
	}

	if adventurer.Actions != "AADADAGGA" {
		t.Errorf("Expected adventurer actions 'AADADAGGA', got '%s'", adventurer.Actions)
	}
}

func TestLoadGameData_EmptyFile(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_empty_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()

	// When
	gameParser := parser.NewParser()
	mapData, err := gameParser.LoadGameData(tempFile.Name())

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(mapData.Mountains) != 0 {
		t.Errorf("Expected 0 mountains, got %d", len(mapData.Mountains))
	}

	if len(mapData.Treasures) != 0 {
		t.Errorf("Expected 0 treasures, got %d", len(mapData.Treasures))
	}

	if len(mapData.Adventurers) != 0 {
		t.Errorf("Expected 0 adventurers, got %d", len(mapData.Adventurers))
	}
}

func TestLoadGameData_WithComments(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_comments_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	gameContent := `# This is a comment
C - 3 - 4
# Another comment
M - 1 - 0

T - 0 - 3 - 2
# Final comment
A - Lara - 1 - 1 - S - AADADAGGA`

	tempFile.WriteString(gameContent)
	tempFile.Close()

	// When
	gameParser := parser.NewParser()
	mapData, err := gameParser.LoadGameData(tempFile.Name())

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if mapData.Width != 3 || mapData.Height != 4 {
		t.Errorf("Expected map dimensions 3x4, got %dx%d", mapData.Width, mapData.Height)
	}

	if len(mapData.Mountains) != 1 {
		t.Errorf("Expected 1 mountain, got %d", len(mapData.Mountains))
	}

	if len(mapData.Treasures) != 1 {
		t.Errorf("Expected 1 treasure, got %d", len(mapData.Treasures))
	}

	if len(mapData.Adventurers) != 1 {
		t.Errorf("Expected 1 adventurer, got %d", len(mapData.Adventurers))
	}
}

func TestLoadGameData_NonExistentFile(t *testing.T) {
	// When
	gameParser := parser.NewParser()
	_, err := gameParser.LoadGameData("nonexistent_file.txt")

	// Then
	if err == nil {
		t.Errorf("Expected error for non-existent file, but got none")
	}
}
