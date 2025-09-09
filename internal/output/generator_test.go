package output_test

import (
	"carte_au_tresor/internal/domain/models"
	"carte_au_tresor/internal/output"
	"os"
	"testing"
)

func TestGenerateOutputLines(t *testing.T) {
	// Given
	gameMapData := &models.MapData{
		Width:  3,
		Height: 4,
		Mountains: []models.Mountain{
			{Row: 0, Col: 1},
			{Row: 1, Col: 2},
		},
		Treasures: []models.Treasure{
			{Row: 3, Col: 0, Count: 2},
			{Row: 3, Col: 1, Count: 3},
		},
		Adventurers: []models.Adventurer{
			{
				Name:          "Lara",
				Row:           3,
				Col:           0,
				Direction:     models.South,
				Actions:       "",
				ActionIndex:   0,
				TreasureCount: 3,
			},
		},
	}

	// Create a game map with treasures and adventurer
	gameMap := make([][]models.MapCell, gameMapData.Height)
	for i := range gameMap {
		gameMap[i] = make([]models.MapCell, gameMapData.Width)
	}

	// Place mountains
	gameMap[0][1].Mountain = &gameMapData.Mountains[0]
	gameMap[1][2].Mountain = &gameMapData.Mountains[1]

	// Place remaining treasure (after Lara collected some)
	treasure := models.Treasure{Row: 3, Col: 1, Count: 2}
	gameMap[3][1].Treasure = &treasure

	// Place adventurer
	gameMap[3][0].Adventurer = &gameMapData.Adventurers[0]

	// When
	generator := output.NewGenerator()
	lines := generator.GenerateOutputLines(gameMapData, &gameMap)

	// Then
	expected := []string{
		"C - 3 - 4",
		"M - 1 - 0",
		"M - 2 - 1",
		"T - 1 - 3 - 2",
		"A - Lara - 0 - 3 - S - 3",
	}

	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
		for i, line := range lines {
			t.Logf("Line %d: %s", i, line)
		}
	}

	for i, expectedLine := range expected {
		if i < len(lines) && lines[i] != expectedLine {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLine, lines[i])
		}
	}
}

func TestGenerateOutputLines_NoTreasuresLeft(t *testing.T) {
	// Given
	gameMapData := &models.MapData{
		Width:     3,
		Height:    4,
		Mountains: []models.Mountain{{Row: 0, Col: 1}},
		Treasures: []models.Treasure{},
		Adventurers: []models.Adventurer{
			{
				Name:          "Lara",
				Row:           1,
				Col:           1,
				Direction:     models.North,
				Actions:       "",
				ActionIndex: 0,
				TreasureCount: 1,
			},
		},
	}

	gameMap := make([][]models.MapCell, gameMapData.Height)
	for i := range gameMap {
		gameMap[i] = make([]models.MapCell, gameMapData.Width)
	}

	gameMap[0][1].Mountain = &gameMapData.Mountains[0]
	gameMap[1][1].Adventurer = &gameMapData.Adventurers[0]

	// When
	generator := output.NewGenerator()
	lines := generator.GenerateOutputLines(gameMapData, &gameMap)

	// Then
	expected := []string{
		"C - 3 - 4",
		"M - 1 - 0",
		"A - Lara - 1 - 1 - N - 1",
	}

	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, expectedLine := range expected {
		if i < len(lines) && lines[i] != expectedLine {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLine, lines[i])
		}
	}
}

func TestGenerateOutputLines_EmptyMap(t *testing.T) {
	// Given
	gameMapData := &models.MapData{
		Width:       2,
		Height:      2,
		Mountains:   []models.Mountain{},
		Treasures:   []models.Treasure{},
		Adventurers: []models.Adventurer{},
	}

	gameMap := make([][]models.MapCell, gameMapData.Height)
	for i := range gameMap {
		gameMap[i] = make([]models.MapCell, gameMapData.Width)
	}

	// When
	generator := output.NewGenerator()
	lines := generator.GenerateOutputLines(gameMapData, &gameMap)

	// Then
	expectedLines := []string{
		"C - 2 - 2",
	}

	if len(lines) != len(expectedLines) {
		t.Errorf("Expected %d lines, got %d", len(expectedLines), len(lines))
	}

	for i, expected := range expectedLines {
		if i < len(lines) && lines[i] != expected {
			t.Errorf("Line %d: expected %q, got %q", i, expected, lines[i])
		}
	}
}

func TestGenerateOutputLines_MultipleAdventurers(t *testing.T) {
	// Given
	gameMapData := &models.MapData{
		Width:     3,
		Height:    3,
		Mountains: []models.Mountain{},
		Treasures: []models.Treasure{},
		Adventurers: []models.Adventurer{
			{
				Name:          "Lara",
				Row:           0,
				Col:           0,
				Direction:     models.East,
				Actions:       "",
				ActionIndex: 0,
				TreasureCount: 1,
			},
			{
				Name:          "Indiana",
				Row:           2,
				Col:           2,
				Direction:     models.West,
				Actions:       "",
				ActionIndex: 0,
				TreasureCount: 0,
			},
		},
	}

	gameMap := make([][]models.MapCell, gameMapData.Height)
	for i := range gameMap {
		gameMap[i] = make([]models.MapCell, gameMapData.Width)
	}

	gameMap[0][0].Adventurer = &gameMapData.Adventurers[0]
	gameMap[2][2].Adventurer = &gameMapData.Adventurers[1]

	// When
	generator := output.NewGenerator()
	lines := generator.GenerateOutputLines(gameMapData, &gameMap)

	// Then
	expected := []string{
		"C - 3 - 3",
		"A - Lara - 0 - 0 - E - 1",
		"A - Indiana - 2 - 2 - O - 0",
	}

	if len(lines) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(lines))
	}

	for i, expectedLine := range expected {
		if i < len(lines) && lines[i] != expectedLine {
			t.Errorf("Line %d: expected %q, got %q", i, expectedLine, lines[i])
		}
	}
}

func TestWriteGameOutput(t *testing.T) {
	// Given
	tempFile, err := os.CreateTemp("", "test_output_*.txt")
	if err != nil {
		t.Fatal(err)
	}
	tempFile.Close()
	defer os.Remove(tempFile.Name())

	gameMapData := &models.MapData{
		Width:  2,
		Height: 2,
		Mountains: []models.Mountain{
			{Row: 0, Col: 1},
		},
		Treasures: []models.Treasure{},
		Adventurers: []models.Adventurer{
			{
				Name:          "Test",
				Row:           1,
				Col:           0,
				Direction:     models.North,
				Actions:       "",
				ActionIndex: 0,
				TreasureCount: 0,
			},
		},
	}

	gameMap := make([][]models.MapCell, gameMapData.Height)
	for i := range gameMap {
		gameMap[i] = make([]models.MapCell, gameMapData.Width)
	}

	gameMap[0][1].Mountain = &gameMapData.Mountains[0]
	gameMap[1][0].Adventurer = &gameMapData.Adventurers[0]

	// When
	generator := output.NewGenerator()
	err = generator.WriteGameOutput(gameMapData, &gameMap, tempFile.Name())

	// Then
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify file content
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expectedContent := `C - 2 - 2
M - 1 - 0
A - Test - 0 - 1 - N - 0
`

	if string(content) != expectedContent {
		t.Errorf("Expected:\n%s\nGot:\n%s", expectedContent, string(content))
	}
}
