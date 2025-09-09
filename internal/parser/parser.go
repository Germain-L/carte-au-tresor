package parser

import (
	"carte_au_tresor/internal/domain/models"
	"carte_au_tresor/internal/fileutil"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Parser struct{}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) parseCoords(line string) (row, col int, err error) {
	coordRe := regexp.MustCompile(`^\D*(\d+)\D+(\d+)`)

	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "#") || line == "" {
		return 0, 0, fmt.Errorf("skipped line: %q", line)
	}

	m := coordRe.FindStringSubmatch(line)
	if m == nil {
		return 0, 0, fmt.Errorf("no coordinates in %q", line)
	}

	col, err = strconv.Atoi(m[1]) // X coordinate (column)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid col %q: %w", m[1], err)
	}

	row, err = strconv.Atoi(m[2]) // Y coordinate (row)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid row %q: %w", m[2], err)
	}

	return row, col, nil
}

func (p *Parser) ParseMapDimensions(line string) (width, height int, err error) {
	// Note: parseCoords returns (row, col) but for map dimensions,
	// the format is "C - width - height", so we need to swap them
	height, width, err = p.parseCoords(line)
	return
}

func (p *Parser) ParseMountain(line string) (models.Mountain, error) {
	row, col, err := p.parseCoords(line)
	if err != nil {
		return models.Mountain{}, err
	}
	return models.Mountain{Row: row, Col: col}, nil
}

func (p *Parser) ParseTreasure(line string) (models.Treasure, error) {
	row, col, err := p.parseCoords(line)
	if err != nil {
		return models.Treasure{}, err
	}

	re := regexp.MustCompile(`(?:\s*-\s*\d+){2}\s*-\s*(\d+)$`)
	m := re.FindStringSubmatch(line)
	if m == nil {
		return models.Treasure{}, fmt.Errorf("invalid treasure line in %q", line)
	}

	nbTreasures, err := strconv.Atoi(m[1])
	if err != nil {
		return models.Treasure{}, err
	}

	return models.Treasure{Row: row, Col: col, Count: nbTreasures}, nil
}

func (p *Parser) ParseAdventurer(line string) (models.Adventurer, error) {
	row, col, err := p.parseCoords(line)
	if err != nil {
		return models.Adventurer{}, err
	}

	adventurerRe := regexp.MustCompile(`^A\s*-\s*([^-]+)\s*-\s*\d+\s*-\s*\d+\s*-\s*([^-]+)\s*-\s*(.+)$`)
	m := adventurerRe.FindStringSubmatch(line)
	if m == nil {
		return models.Adventurer{}, fmt.Errorf("invalid adventurer line format in %q", line)
	}

	name := strings.TrimSpace(m[1])
	directionStr := strings.TrimSpace(m[2])
	actions := strings.TrimSpace(m[3])

	// Convert string direction to Direction constant
	direction, exists := models.StringToDirection[directionStr]
	if !exists {
		return models.Adventurer{}, fmt.Errorf("invalid direction %q", directionStr)
	}

	return models.Adventurer{
		Name:          name,
		Row:           row,
		Col:           col,
		Direction:     direction,
		Actions:       actions,
		TreasureCount: 0,
	}, nil
}

func (p *Parser) LoadGameData(filename string) (*models.MapData, error) {
	lines, errorCh := fileutil.ReadToChannel(filename)

	gameMapData := &models.MapData{
		Mountains:   []models.Mountain{},
		Treasures:   []models.Treasure{},
		Adventurers: []models.Adventurer{},
	}

	for line := range lines {
		if len(line) <= 0 {
			continue
		}

		switch line[0] {
		case 'C':
			width, height, err := p.ParseMapDimensions(line)
			if err != nil {
				return nil, err
			}
			gameMapData.Width = width
			gameMapData.Height = height

		case 'M':
			mountain, err := p.ParseMountain(line)
			if err != nil {
				return nil, err
			}
			gameMapData.Mountains = append(gameMapData.Mountains, mountain)

		case 'T':
			treasure, err := p.ParseTreasure(line)
			if err != nil {
				return nil, err
			}
			gameMapData.Treasures = append(gameMapData.Treasures, treasure)

		case 'A':
			adventurer, err := p.ParseAdventurer(line)
			if err != nil {
				return nil, err
			}
			gameMapData.Adventurers = append(gameMapData.Adventurers, adventurer)
		}
	}

	// Check for file reading errors
	err := <-errorCh
	if err != nil {
		return nil, err
	}

	return gameMapData, nil
}
