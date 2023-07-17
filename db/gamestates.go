package db

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/quevivasbien/bird-backend/game"
)

const GAME_TABLE_NAME = "GameTable"

type GameTable struct {
	client *dynamodb.Client
}

// initialize a new GameTable struct with given client
// if table already exists, use that table
// otherwise, create the table
func MakeGameTable(client *dynamodb.Client) (GameTable, error) {
	exists, err := TableExists(client, GAME_TABLE_NAME)
	if err != nil {
		return GameTable{}, fmt.Errorf("Error when checking if game table exists: %v", err)
	}
	if exists {
		return GameTable{client}, nil
	} else {
		return createGameTable(client)
	}
}

// remove the table from dynamodb, useful for a complete reset
func DeleteGameTable(client *dynamodb.Client) error {
	_, err := client.DeleteTable(
		context.TODO(),
		&dynamodb.DeleteTableInput{
			TableName: aws.String(GAME_TABLE_NAME),
		},
	)
	if err != nil {
		return fmt.Errorf("Error when deleting game table: %v", err)
	}
	return nil
}

func createGameTable(client *dynamodb.Client) (GameTable, error) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("GameID"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("GameID"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(GAME_TABLE_NAME),
		BillingMode: types.BillingModePayPerRequest,
	}
	_, err := client.CreateTable(context.TODO(), input)
	if err != nil {
		return GameTable{}, fmt.Errorf("Error when creating game table: %v", err)
	}
	return GameTable{client}, nil
}

func (t GameTable) GetGameState(id string) (game.GameState, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"GameID": &types.AttributeValueMemberS{Value: id},
		},
		TableName: aws.String(GAME_TABLE_NAME),
	}
	output, err := t.client.GetItem(context.TODO(), input)
	if err != nil {
		return game.GameState{}, fmt.Errorf("Error when fetching game state: %v", err)
	}
	if output.Item == nil {
		return game.GameState{}, fmt.Errorf("No game found with ID %v", id)
	}
	gameState := game.GameState{}
	err = attributevalue.UnmarshalMap(output.Item, &gameState)
	if err != nil {
		return gameState, fmt.Errorf("Error when unpacking game state: %v", err)
	}
	return gameState, nil
}

func (t GameTable) PutGameState(s game.GameState) error {
	item, err := attributevalue.MarshalMap(s)
	if err != nil {
		return fmt.Errorf("Error when packing game state: %v", err)
	}
	_, err = t.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(GAME_TABLE_NAME),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("Error adding game state to database: %v", err)
	}
	return nil
}

func (t GameTable) UpdateGameState(gameID string, updates map[string]interface{}) error {
	update := expression.UpdateBuilder{}
	for key, value := range updates {
		update = update.Set(expression.Name(key), expression.Value(value))
	}
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("Error when building update expression: %v", err)
	}
	_, err = t.client.UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName: aws.String(GAME_TABLE_NAME),
			Key: map[string]types.AttributeValue{
				"GameID": &types.AttributeValueMemberS{Value: gameID},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		},
	)
	if err != nil {
		return fmt.Errorf("Error when updating game state: %v", err)
	}
	return nil
}
