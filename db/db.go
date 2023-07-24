package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const BILLING_MODE = types.BillingModePayPerRequest

func GetClient(region string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}

type Table interface {
	Client() *dynamodb.Client
	Name() string
	IndexName() string
	IndexType() types.ScalarAttributeType
}

// remove the table from dynamodb, useful for a complete reset
func deleteTable(t Table) error {
	_, err := t.Client().DeleteTable(
		context.TODO(),
		&dynamodb.DeleteTableInput{
			TableName: aws.String(t.Name()),
		},
	)
	if err != nil {
		return fmt.Errorf("Error when deleting game table: %v", err)
	}
	return nil
}

func initTable(t Table) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(t.Name()),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String(t.IndexName()),
				AttributeType: t.IndexType(),
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String(t.IndexName()),
				KeyType:       types.KeyTypeHash,
			},
		},
		BillingMode: BILLING_MODE,
	}
	_, err := t.Client().CreateTable(context.TODO(), input)
	if err != nil {
		return fmt.Errorf("Error when creating table %s: %v", t.Name(), err)
	}
	return nil
}

func tableIsInitialized(t Table) (bool, error) {
	_, err := t.Client().DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{TableName: aws.String(t.Name())},
	)
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

type ItemNotFound struct {
	ItemName string
}

func (i ItemNotFound) Error() string {
	if i.ItemName == "" {
		return "Item not found in database"
	} else {
		return fmt.Sprintf("%s not found in database", i.ItemName)
	}
}

func getItem(t Table, id string) (map[string]types.AttributeValue, error) {
	input := &dynamodb.GetItemInput{
		TableName: aws.String(t.Name()),
		Key: map[string]types.AttributeValue{
			t.IndexName(): &types.AttributeValueMemberS{Value: id},
		},
	}
	output, err := t.Client().GetItem(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("Error when fetching item %s from table %s: %v", id, t.Name(), err)
	}
	return output.Item, nil
}

func putItem(t Table, item interface{}) error {
	itemMap, err := attributevalue.MarshalMap(item)
	if err != nil {
		return fmt.Errorf("Error when packing item to be placed in table %s: %v", t.Name(), err)
	}
	_, err = t.Client().PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(t.Name()),
		Item:      itemMap,
	})
	if err != nil {
		return fmt.Errorf("Error adding item to table %s: %v", t.Name(), err)
	}
	return nil
}

func updateItem(t Table, id string, updates map[string]interface{}) error {
	update := expression.UpdateBuilder{}
	for key, value := range updates {
		update = update.Set(expression.Name(key), expression.Value(value))
	}
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return fmt.Errorf("Error when building update expression: %v", err)
	}
	_, err = t.Client().UpdateItem(
		context.TODO(),
		&dynamodb.UpdateItemInput{
			TableName: aws.String(t.Name()),
			Key: map[string]types.AttributeValue{
				t.IndexName(): &types.AttributeValueMemberS{Value: id},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		},
	)
	if err != nil {
		return fmt.Errorf("Error when updating item in %s on db: %v", t.Name(), err)
	}
	return nil
}

func deleteItem(t Table, id string) error {
	_, err := t.Client().DeleteItem(
		context.TODO(),
		&dynamodb.DeleteItemInput{
			TableName: aws.String(t.Name()),
			Key: map[string]types.AttributeValue{
				t.IndexName(): &types.AttributeValueMemberS{Value: id},
			},
		},
	)
	if err != nil {
		return fmt.Errorf("Error when deleting item %s from table %s: %v", id, t.Name(), err)
	}
	return nil
}

type Tables struct {
	Region string
	UserTable
	LobbyTable
	BidTable
	GameTable
}

func GetTables(region string) (Tables, error) {
	client, err := GetClient(region)
	if err != nil {
		return Tables{}, fmt.Errorf("Error getting database client: %v", err)
	}
	tables := Tables{Region: region}
	tables.UserTable, err = MakeUserTable(client)
	if err != nil {
		return tables, fmt.Errorf("Error initializing user table: %v", err)
	}
	tables.LobbyTable, err = MakeLobbyTable(client)
	if err != nil {
		return tables, fmt.Errorf("Error initializing lobby table: %v", err)
	}
	tables.BidTable, err = MakeBidTable(client)
	if err != nil {
		return tables, fmt.Errorf("Error initializing bid table: %v", err)
	}
	tables.GameTable, err = MakeGameTable(client)
	if err != nil {
		return tables, fmt.Errorf("Error initializing game table: %v", err)
	}
	return tables, nil
}

// delete and re-initialize all tables
func (t *Tables) Reset() error {
	err := deleteTable(t.BidTable)
	if err != nil {
		return fmt.Errorf("Problem deleting bid table: %v", err)
	}
	err = deleteTable(t.GameTable)
	if err != nil {
		return fmt.Errorf("Problem deleting game table: %v", err)
	}
	err = deleteTable(t.LobbyTable)
	if err != nil {
		return fmt.Errorf("Problem deleting lobby table: %v", err)
	}
	err = deleteTable(t.UserTable)
	if err != nil {
		return fmt.Errorf("Problem deleting user table: %v", err)
	}
	*t, err = GetTables(t.Region)
	if err != nil {
		return fmt.Errorf("Problem re-initializing tables: %v", err)
	}
	return nil
}
