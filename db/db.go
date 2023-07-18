package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func GetClient(region string) (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}
	return dynamodb.NewFromConfig(cfg), nil
}

func TableExists(client *dynamodb.Client, name string) (bool, error) {
	_, err := client.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{TableName: aws.String(name)},
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

type Tables struct {
	UserTable
	GameTable
}

func GetTables(region string) (Tables, error) {
	client, err := GetClient(region)
	if err != nil {
		return Tables{}, fmt.Errorf("Error getting database client: %v", err)
	}
	tables := Tables{}
	tables.UserTable, err = MakeUserTable(client)
	if err != nil {
		return tables, fmt.Errorf("Error initializing user table: %v", err)
	}
	tables.GameTable, err = MakeGameTable(client)
	if err != nil {
		return tables, fmt.Errorf("Error initializing game table: %v", err)
	}
	return tables, nil
}
