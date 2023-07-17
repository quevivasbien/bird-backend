package main

import (
	"fmt"

	"github.com/quevivasbien/bird-backend/db"
	"github.com/quevivasbien/bird-backend/game"
)

func main() {
	client, err := db.GetClient("us-east-1")
	if err != nil {
		panic(fmt.Sprintf("Error when fetching client: %v", err))
	}
	gameTable, err := db.MakeGameTable(client)
	if err != nil {
		panic(fmt.Sprintf("Error when getting game table: %v", err))
	}

	gameState := game.GameState{
		GameID: "TestGame",
	}
	// err = gameTable.PutGameState(gameState)
	// if err != nil {
	// 	panic(fmt.Sprintf("Error when putting game state in db: %v", err))
	// }

	gameState.Players[0] = "bob"
	gameState.Bid = 150
	err = gameTable.UpdateGameState(gameState)
	if err != nil {
		panic(fmt.Sprintf("Error when updating game state on db: %v", err))
	}

	response, err := gameTable.GetGameState(gameState.GameID)
	if err != nil {
		panic(fmt.Sprintf("Error when retrieving game state from db: %v", err))
	}
	fmt.Printf("Got game state %v\n", response)
}
