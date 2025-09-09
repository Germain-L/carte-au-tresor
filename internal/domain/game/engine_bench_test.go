package game

import (
	"carte_au_tresor/internal/domain/models"
	"testing"
)

func BenchmarkDirectionOperations(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = TurnLeft(models.North)
		_ = TurnLeft(models.South)
		_ = TurnLeft(models.East)
		_ = TurnLeft(models.West)
		_ = TurnRight(models.North)
		_ = TurnRight(models.South)
		_ = TurnRight(models.East)
		_ = TurnRight(models.West)
	}
}

func BenchmarkGameSimulation(b *testing.B) {
	// Create a sample map for benchmarking
	mapData := &models.MapData{
		Width:  10,
		Height: 10,
		Mountains: []models.Mountain{
			{Row: 1, Col: 1},
			{Row: 5, Col: 5},
		},
		Treasures: []models.Treasure{
			{Row: 3, Col: 3, Count: 5},
			{Row: 7, Col: 7, Count: 3},
		},
		Adventurers: []models.Adventurer{
			{
				Name:          "TestAdventurer",
				Row:           0,
				Col:           0,
				Direction:     models.South,
				Actions:       "AADADAGGGAAADDDAAAGGGAAADDDAAA",
				TreasureCount: 0,
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create fresh copy for each benchmark run
		engine := NewEngine(mapData)
		engine.Run()
	}
}

func BenchmarkActionProcessing(b *testing.B) {
	mapData := &models.MapData{
		Width:  5,
		Height: 5,
		Adventurers: []models.Adventurer{
			{
				Name:          "TestAdventurer",
				Row:           2,
				Col:           2,
				Direction:     models.North,
				Actions:       "",
				TreasureCount: 0,
			},
		},
	}
	
	engine := NewEngine(mapData)
	adventurer := &engine.mapData.Adventurers[0]

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.ProcessAction(adventurer, 'G')
		engine.ProcessAction(adventurer, 'D')
		engine.ProcessAction(adventurer, 'A')
	}
}