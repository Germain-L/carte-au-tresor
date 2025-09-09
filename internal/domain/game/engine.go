package game

import (
	"carte_au_tresor/internal/domain/models"
	"fmt"
)

type Engine struct {
	mapData *models.MapData
	gameMap [][]models.MapCell
}

func NewEngine(mapData *models.MapData) *Engine {
	engine := &Engine{
		mapData: mapData,
	}
	engine.gameMap = engine.populateGameMap()
	return engine
}

func (e *Engine) GetMapData() *models.MapData {
	return e.mapData
}

func (e *Engine) GetGameMap() *[][]models.MapCell {
	return &e.gameMap
}

func (e *Engine) populateGameMap() [][]models.MapCell {
	gameMap := make([][]models.MapCell, e.mapData.Height)
	for i := range gameMap {
		gameMap[i] = make([]models.MapCell, e.mapData.Width)
	}

	// Insert mountains
	for _, mountain := range e.mapData.Mountains {
		gameMap[mountain.Row][mountain.Col].Mountain = &mountain
	}

	// Insert treasures
	for _, treasure := range e.mapData.Treasures {
		gameMap[treasure.Row][treasure.Col].Treasure = &treasure
	}

	// Insert adventurers
	for _, adventurer := range e.mapData.Adventurers {
		gameMap[adventurer.Row][adventurer.Col].Adventurer = &adventurer
	}

	return gameMap
}

func (e *Engine) CanMove(nextRow, nextCol int) bool {
	// Check boundaries, not defined in pdf, so assuming players don't wrap around the map
	if nextRow < 0 || nextRow >= len(e.gameMap) || nextCol < 0 || nextCol >= len(e.gameMap[0]) {
		return false
	}

	cell := e.gameMap[nextRow][nextCol]
	return cell.Mountain == nil && cell.Adventurer == nil
}

func TurnLeft(direction models.Direction) models.Direction {
	switch direction {
	case models.North:
		return models.West
	case models.South:
		return models.East
	case models.East:
		return models.North
	case models.West:
		return models.South
	default:
		return direction
	}
}

func TurnRight(direction models.Direction) models.Direction {
	switch direction {
	case models.North:
		return models.East
	case models.South:
		return models.West
	case models.East:
		return models.South
	case models.West:
		return models.North
	default:
		return direction
	}
}

func (e *Engine) MoveAdventurer(adventurer *models.Adventurer) {
	var nextRow, nextCol int

	switch adventurer.Direction {
	case models.East:
		nextRow = adventurer.Row
		nextCol = adventurer.Col + 1
	case models.West:
		nextRow = adventurer.Row
		nextCol = adventurer.Col - 1
	case models.South:
		nextRow = adventurer.Row + 1
		nextCol = adventurer.Col
	case models.North:
		nextRow = adventurer.Row - 1
		nextCol = adventurer.Col
	default:
		return
	}

	if e.CanMove(nextRow, nextCol) {
		// Remove adventurer from old position
		e.gameMap[adventurer.Row][adventurer.Col].Adventurer = nil

		// Update adventurer position
		adventurer.Row = nextRow
		adventurer.Col = nextCol

		// Check if there's a treasure at the new position
		cell := &e.gameMap[nextRow][nextCol]
		if cell.Treasure != nil && cell.Treasure.Count > 0 {
			// Collect one treasure
			cell.Treasure.Count--
			adventurer.TreasureCount++  // Optimized: simple counter instead of slice
			fmt.Printf("Adventurer %s collected a treasure at (%d, %d). Remaining treasures: %d\n",
				adventurer.Name, nextRow, nextCol, cell.Treasure.Count)

			// If no more treasures, remove the treasure from the cell
			if cell.Treasure.Count == 0 {
				cell.Treasure = nil
			}
		}

		// Place adventurer at new position
		cell.Adventurer = adventurer
	} else {
		fmt.Printf("Adventurer %s cannot move to (%d, %d)\n", adventurer.Name, nextRow, nextCol)
	}
}

func (e *Engine) ProcessAction(adventurer *models.Adventurer, action rune) {
	switch action {
	case 'G':
		adventurer.Direction = TurnLeft(adventurer.Direction)
	case 'D':
		adventurer.Direction = TurnRight(adventurer.Direction)
	case 'A':
		e.MoveAdventurer(adventurer)
	default:
		// Do nothing for letters we don't recognize
	}
}

func (e *Engine) Run() {
	for {
		// Default to true to not loop again if no actions are left
		allActionsCompleted := true

		// Process each adventuer in order
		for i := range e.mapData.Adventurers {
			if len(e.mapData.Adventurers[i].Actions) > 0 {
				allActionsCompleted = false

				// Convert to rune (Go's type for single characters)
				nextAction := rune(e.mapData.Adventurers[i].Actions[0])
				e.mapData.Adventurers[i].Actions = e.mapData.Adventurers[i].Actions[1:]

				e.ProcessAction(&e.mapData.Adventurers[i], nextAction)
			}
		}

		if allActionsCompleted {
			break
		}
	}
}
