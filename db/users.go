package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

const USER_TABLE_NAME = "BirdUsers"

type UserTable struct {
	client *dynamodb.Client
}

// initialize a new GameTable struct with given client
// if table already exists, use that table
// otherwise, create the table
func MakeUserTable(client *dynamodb.Client) (UserTable, error) {
	exists, err := TableExists(client, USER_TABLE_NAME)
	if err != nil {
		return UserTable{}, fmt.Errorf("Error when checking if game table exists: %v", err)
	}
	if exists {
		return UserTable{client}, nil
	} else {
		return createUserTable(client)
	}
}

// remove the table from dynamodb, useful for a complete reset
func DeleteUserTable(client *dynamodb.Client) error {
	_, err := client.DeleteTable(
		context.TODO(),
		&dynamodb.DeleteTableInput{
			TableName: aws.String(USER_TABLE_NAME),
		},
	)
	if err != nil {
		return fmt.Errorf("Error when deleting game table: %v", err)
	}
	return nil
}

func createUserTable(client *dynamodb.Client) (UserTable, error) {
	input := &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Name"),
				KeyType:       types.KeyTypeHash,
			},
		},
		TableName:   aws.String(USER_TABLE_NAME),
		BillingMode: types.BillingModePayPerRequest,
	}
	_, err := client.CreateTable(context.TODO(), input)
	if err != nil {
		return UserTable{}, fmt.Errorf("Error when creating user table: %v", err)
	}
	return UserTable{client}, nil
}

func (t UserTable) GetUser(uname string) (User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"Name": &types.AttributeValueMemberS{Value: uname},
		},
		TableName: aws.String(USER_TABLE_NAME),
	}
	output, err := t.client.GetItem(context.TODO(), input)
	if err != nil {
		return User{}, fmt.Errorf("Error when fetching user: %v", err)
	}
	if output.Item == nil {
		return User{}, fmt.Errorf("No user found with name %v", uname)
	}
	user := User{}
	err = attributevalue.UnmarshalMap(output.Item, &user)
	if err != nil {
		return user, fmt.Errorf("Error when unpacking user: %v", err)
	}
	return user, nil
}

func (t UserTable) PutUser(u User) error {
	item, err := attributevalue.MarshalMap(u)
	if err != nil {
		return fmt.Errorf("Error when packing user: %v", err)
	}
	_, err = t.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(USER_TABLE_NAME),
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("Error adding user to database: %v", err)
	}
	return nil
}

func (t UserTable) UpdateUser(uname string, updates map[string]interface{}) error {
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
			TableName: aws.String(USER_TABLE_NAME),
			Key: map[string]types.AttributeValue{
				"Name": &types.AttributeValueMemberS{Value: uname},
			},
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
		},
	)
	if err != nil {
		return fmt.Errorf("Error when updating user in db: %v", err)
	}
	return nil
}

// check that user exists and has correct password
func (t UserTable) ValidateUser(name string, password string) (bool, User, error) {
	dbUser, err := t.GetUser(name)
	if err != nil {
		var notFound *types.ResourceNotFoundException
		if errors.As(err, &notFound) {
			return false, User{}, nil
		}
		return false, User{}, err
	}
	ok := dbUser.Password == password
	var user User
	if ok {
		user = dbUser
	} else {
		user = User{}
	}
	return ok, user, nil
}
