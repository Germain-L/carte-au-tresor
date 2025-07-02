package game_test

import (
	"carte_au_tresor/internal/domain/game"
	"carte_au_tresor/internal/domain/models"
	"testing"
)

func TestTurnLeft(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"N", "O"},
		{"S", "E"},
		{"E", "N"},
		{"O", "S"},
		{"X", "X"}, // Invalid direction should remain unchanged
	}

	for _, test := range tests {
		result := game.TurnLeft(test.input)
		if result != test.expected {
			t.Errorf("TurnLeft(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestTurnRight(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"N", "E"},
		{"S", "O"},
		{"E", "S"},
		{"O", "N"},
		{"X", "X"}, // Invalid direction should remain unchanged
	}

	for _, test := range tests {
		result := game.TurnRight(test.input)
		if result != test.expected {
			t.Errorf("TurnRight(%s): expected %s, got %s", test.input, test.expected, result)
		}
	}
}

func TestCanMove(t *testing.T) {
	// Given
	mapData := &models.MapData{
		Width:  3,
		Height: 3,
		Mountains: []models.Mountain{
			{Row: 1, Col: 1},
		},
		Adventurers: []models.Adventurer{
			{Row: 2, Col: 2},
		},
	}

	engine := game.NewEngine(mapData)

	tests := []struct {
		row      int
		col      int
		expected bool
	}{
		{0, 0, true},   // Empty cell
		{1, 1, false},  // Mountain
		{2, 2, false},  // Adventurer
		{-1, 0, false}, // Out of bounds
		{0, -1, false}, // Out of bounds
		{3, 0, false},  // Out of bounds
		{0, 3, false},  // Out of bounds
	}

	for _, test := range tests {
		result := engine.CanMove(test.row, test.col)
		if result != test.expected {
			t.Errorf("CanMove(%d, %d): expected %t, got %t", test.row, test.col, test.expected, result)
		}
	}
}

func TestProcessAction(t *testing.T) {
	// Given
	mapData := &models.MapData{
		Width:  3,
		Height: 3,
		Adventurers: []models.Adventurer{
			{
				Name:      "Test",
				Row:       1,
				Col:       1,
				Direction: "N",
				Actions:   "",
				Treasures: []models.Treasure{},
			},
		},
	}

	engine := game.NewEngine(mapData)
	adventurer := &mapData.Adventurers[0]

	// Test turn left
	engine.ProcessAction(adventurer, 'G')
	if adventurer.Direction != "O" {
		t.Errorf("Expected direction O after turning left from N, got %s", adventurer.Direction)
	}

	// Test turn right
	engine.ProcessAction(adventurer, 'D')
	if adventurer.Direction != "N" {
		t.Errorf("Expected direction N after turning right from O, got %s", adventurer.Direction)
	}

	// Test move (should move north from position 1,1 to 0,1)
	engine.ProcessAction(adventurer, 'A')
	if adventurer.Row != 0 || adventurer.Col != 1 {
		t.Errorf("Expected position (0,1) after moving north, got (%d,%d)", adventurer.Row, adventurer.Col)
	}
}

func TestMoveAdventurer_TreasureCollection(t *testing.T) {
	// Given
	treasure := models.Treasure{Row: 1, Col: 1, Count: 2}
	mapData := &models.MapData{
		Width:     3,
		Height:    3,
		Treasures: []models.Treasure{treasure},
		Adventurers: []models.Adventurer{
			{
				Name:      "Lara",
				Row:       0,
				Col:       1,
				Direction: "S",
				Actions:   "",
				Treasures: []models.Treasure{},
			},
		},
	}

	engine := game.NewEngine(mapData)
	adventurer := &mapData.Adventurers[0]

	// When - move south to collect treasure
	engine.MoveAdventurer(adventurer)

	// Then
	if adventurer.Row != 1 || adventurer.Col != 1 {
		t.Errorf("Expected adventurer at (1,1), got (%d,%d)", adventurer.Row, adventurer.Col)
	}

	if len(adventurer.Treasures) != 1 {
		t.Errorf("Expected adventurer to have 1 treasure, got %d", len(adventurer.Treasures))
	}

	// Check that treasure count decreased on the map
	gameMap := engine.GetGameMap()
	cell := (*gameMap)[1][1]
	if cell.Treasure == nil || cell.Treasure.Count != 1 {
		t.Errorf("Expected 1 treasure remaining on the map, got %v", cell.Treasure)
	}
}

func TestMoveAdventurer_TreasureExhaustion(t *testing.T) {
	// Given
	treasure := models.Treasure{Row: 1, Col: 1, Count: 1}
	mapData := &models.MapData{
		Width:     3,
		Height:    3,
		Treasures: []models.Treasure{treasure},
		Adventurers: []models.Adventurer{
			{
				Name:      "Lara",
				Row:       0,
				Col:       1,
				Direction: "S",
				Actions:   "",
				Treasures: []models.Treasure{},
			},
		},
	}

	engine := game.NewEngine(mapData)
	adventurer := &mapData.Adventurers[0]

	// When - move south to collect the last treasure
	engine.MoveAdventurer(adventurer)

	// Then
	gameMap := engine.GetGameMap()
	cell := (*gameMap)[1][1]
	if cell.Treasure != nil {
		t.Errorf("Expected no treasure remaining on the map, got %v", cell.Treasure)
	}

	if len(adventurer.Treasures) != 1 {
		t.Errorf("Expected adventurer to have 1 treasure, got %d", len(adventurer.Treasures))
	}
}

func TestMoveAdventurer_BlockedByMountain(t *testing.T) {
	// Given
	mapData := &models.MapData{
		Width:  3,
		Height: 3,
		Mountains: []models.Mountain{
			{Row: 0, Col: 1},
		},
		Adventurers: []models.Adventurer{
			{
				Name:      "Lara",
				Row:       1,
				Col:       1,
				Direction: "N",
				Actions:   "",
				Treasures: []models.Treasure{},
			},
		},
	}

	engine := game.NewEngine(mapData)
	adventurer := &mapData.Adventurers[0]

	// When - try to move north into mountain
	engine.MoveAdventurer(adventurer)

	// Then - adventurer should not move
	if adventurer.Row != 1 || adventurer.Col != 1 {
		t.Errorf("Expected adventurer to stay at (1,1), got (%d,%d)", adventurer.Row, adventurer.Col)
	}
}

func TestMoveAdventurer_BlockedByAnotherAdventurer(t *testing.T) {
	// Given
	mapData := &models.MapData{
		Width:  3,
		Height: 3,
		Adventurers: []models.Adventurer{
			{
				Name:      "Lara",
				Row:       1,
				Col:       1,
				Direction: "N",
				Actions:   "",
				Treasures: []models.Treasure{},
			},
			{
				Name:      "Indiana",
				Row:       0,
				Col:       1,
				Direction: "S",
				Actions:   "",
				Treasures: []models.Treasure{},
			},
		},
	}

	engine := game.NewEngine(mapData)
	adventurer := &mapData.Adventurers[0]

	// When - try to move north into another adventurer
	engine.MoveAdventurer(adventurer)

	// Then - adventurer should not move
	if adventurer.Row != 1 || adventurer.Col != 1 {
		t.Errorf("Expected adventurer to stay at (1,1), got (%d,%d)", adventurer.Row, adventurer.Col)
	}
}

func TestMoveAdventurer_OutOfBounds(t *testing.T) {
	// Given
	mapData := &models.MapData{
		Width:  2,
		Height: 2,
		Adventurers: []models.Adventurer{
			{
				Name:      "Lara",
				Row:       0,
				Col:       0,
				Direction: "N",
				Actions:   "",
				Treasures: []models.Treasure{},
			},
		},
	}

	engine := game.NewEngine(mapData)
	adventurer := &mapData.Adventurers[0]

	// When - try to move north out of bounds
	engine.MoveAdventurer(adventurer)

	// Then - adventurer should not move
	if adventurer.Row != 0 || adventurer.Col != 0 {
		t.Errorf("Expected adventurer to stay at (0,0), got (%d,%d)", adventurer.Row, adventurer.Col)
	}
}

func TestRunSimulation(t *testing.T) {
	// Given
	mapData := &models.MapData{
		Width:  3,
		Height: 3,
		Treasures: []models.Treasure{
			{Row: 1, Col: 1, Count: 1},
		},
		Adventurers: []models.Adventurer{
			{
				Name:      "Lara",
				Row:       0,
				Col:       0,
				Direction: "E",
				Actions:   "ADA", // Move east, turn right, move south to treasure
				Treasures: []models.Treasure{},
			},
		},
	}

	engine := game.NewEngine(mapData)

	// When
	engine.Run()

	// Then
	adventurer := &mapData.Adventurers[0]
	if adventurer.Row != 1 || adventurer.Col != 1 {
		t.Errorf("Expected adventurer at (1,1), got (%d,%d)", adventurer.Row, adventurer.Col)
	}

	if len(adventurer.Treasures) != 1 {
		t.Errorf("Expected adventurer to have 1 treasure, got %d", len(adventurer.Treasures))
	}

	if len(adventurer.Actions) != 0 {
		t.Errorf("Expected all actions to be consumed, got %q", adventurer.Actions)
	}
}

func TestNewEngine(t *testing.T) {
	// Given
	mapData := &models.MapData{
		Width:  2,
		Height: 2,
		Mountains: []models.Mountain{
			{Row: 0, Col: 1},
		},
		Treasures: []models.Treasure{
			{Row: 1, Col: 0, Count: 2},
		},
		Adventurers: []models.Adventurer{
			{Row: 1, Col: 1},
		},
	}

	// When
	engine := game.NewEngine(mapData)

	// Then
	if engine.GetMapData() != mapData {
		t.Errorf("Expected engine to have the provided map data")
	}

	gameMap := engine.GetGameMap()
	if len(*gameMap) != 2 || len((*gameMap)[0]) != 2 {
		t.Errorf("Expected 2x2 game map, got %dx%d", len(*gameMap), len((*gameMap)[0]))
	}

	// Check that entities are properly placed
	if (*gameMap)[0][1].Mountain == nil {
		t.Errorf("Expected mountain at (0,1)")
	}

	if (*gameMap)[1][0].Treasure == nil {
		t.Errorf("Expected treasure at (1,0)")
	}

	if (*gameMap)[1][1].Adventurer == nil {
		t.Errorf("Expected adventurer at (1,1)")
	}
}
