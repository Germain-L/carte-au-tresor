package main

import (
	"carte_au_tresor/internal/domain/game"
	"carte_au_tresor/internal/output"
	"carte_au_tresor/internal/parser"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Welcome, Carte au Trésor!")

	// Initialize components
	gameParser := parser.NewParser()
	gameEngine := &game.Engine{}
	outputGenerator := output.NewGenerator()

	// Load game data
	gameMapData, err := gameParser.LoadGameData("inputs/example.txt")
	if err != nil {
		log.Fatalln("Error loading game data:", err)
	}

	// Initialize game engine
	gameEngine = game.NewEngine(gameMapData)

	// Run the simulation
	gameEngine.Run()

	// Generate and write output
	finalMapData := gameEngine.GetMapData()
	finalGameMap := gameEngine.GetGameMap()

	err = outputGenerator.WriteGameOutput(finalMapData, finalGameMap, "outputs/example.txt")
	if err != nil {
		log.Printf("Error writing output file: %v", err)
	} else {
		fmt.Println("Output written to outputs/example.txt")
	}
}
