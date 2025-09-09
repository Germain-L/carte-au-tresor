package output

import (
	"carte_au_tresor/internal/domain/models"
	"carte_au_tresor/internal/fileutil"
	"fmt"
)

type Generator struct{}

func NewGenerator() *Generator {
	return &Generator{}
}

func (g *Generator) GenerateOutputLines(gameMapData *models.MapData, gameMap *[][]models.MapCell) []string {
	var lines []string

	// Add map dimensions
	lines = append(lines, fmt.Sprintf("C - %d - %d", gameMapData.Width, gameMapData.Height))

	// Add mountains
	for _, mountain := range gameMapData.Mountains {
		lines = append(lines, fmt.Sprintf("M - %d - %d", mountain.Col, mountain.Row))
	}

	// Add treasures that remain
	for i := 0; i < len(*gameMap); i++ {
		for j := 0; j < len((*gameMap)[i]); j++ {
			cell := (*gameMap)[i][j]
			if cell.Treasure != nil && cell.Treasure.Count > 0 {
				lines = append(lines, fmt.Sprintf("T - %d - %d - %d", j, i, cell.Treasure.Count))
			}
		}
	}

	// Add adventurers with their final positions and treasure count
	for _, adventurer := range gameMapData.Adventurers {
		lines = append(lines, fmt.Sprintf("A - %s - %d - %d - %s - %d",
			adventurer.Name, adventurer.Col, adventurer.Row, adventurer.Direction.String(), adventurer.TreasureCount))
	}

	return lines
}

func (g *Generator) WriteGameOutput(gameMapData *models.MapData, gameMap *[][]models.MapCell, filename string) error {
	lines := g.GenerateOutputLines(gameMapData, gameMap)
	return fileutil.WriteLinesToFile(filename, lines)
}
