package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/quevivasbien/bird-backend/game"
)

const GAME_TABLE_NAME = "Bird.Games"

type GameTable struct {
	client *dynamodb.Client
}

func (t GameTable) Client() *dynamodb.Client {
	return t.client
}

func (t GameTable) Name() string {
	return GAME_TABLE_NAME
}

func (t GameTable) IndexName() string {
	return "GameID"
}

func (t GameTable) IndexType() types.ScalarAttributeType {
	return types.ScalarAttributeTypeS
}

// initialize a new GameTable struct with given client
// if table already exists, use that table
// otherwise, create the table
func MakeGameTable(client *dynamodb.Client) (GameTable, error) {
	table := GameTable{client}
	exists, err := tableIsInitialized(table)
	if err != nil {
		return table, fmt.Errorf("Error when checking if game table exists: %v", err)
	}
	if exists {
		return table, nil
	} else {
		err = initTable(table)
		return table, err
	}
}

func (t GameTable) GetGameState(id string) (game.GameState, error) {
	itemMap, err := getItem(t, id)
	if err != nil {
		return game.GameState{}, err
	}
	if itemMap == nil {
		return game.GameState{}, ItemNotFound{"Game"}
	}
	gameState := game.GameState{}
	err = attributevalue.UnmarshalMap(itemMap, &gameState)
	if err != nil {
		return gameState, fmt.Errorf("Error when unpacking game state: %v", err)
	}
	return gameState, nil
}

func (t GameTable) PutGameState(s game.GameState) error {
	return putItem(t, s)
}

func (t GameTable) UpdateGameState(gameID string, updates map[string]interface{}) error {
	return updateItem(t, gameID, updates)
}
