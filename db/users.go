package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	Admin    bool   `json:"admin"`
}

const USER_TABLE_NAME = "Bird.Users"

type UserTable struct {
	client *dynamodb.Client
}

func (t UserTable) Client() *dynamodb.Client {
	return t.client
}

func (t UserTable) Name() string {
	return USER_TABLE_NAME
}

func (t UserTable) IndexName() string {
	return "Name"
}

func (t UserTable) IndexType() types.ScalarAttributeType {
	return types.ScalarAttributeTypeS
}

// initialize a new GameTable struct with given client
// if table already exists, use that table
// otherwise, create the table
func MakeUserTable(client *dynamodb.Client) (UserTable, error) {
	table := UserTable{client}
	exists, err := tableIsInitialized(table)
	if err != nil {
		return table, fmt.Errorf("Error when checking if user table exists: %v", err)
	}
	if exists {
		return table, nil
	} else {
		err = initTable(table)
		return table, err
	}
}

func (t UserTable) GetUser(uname string) (User, error) {
	itemMap, err := getItem(t, uname)
	if err != nil {
		return User{}, err
	}
	if itemMap == nil {
		return User{}, ItemNotFound{"User"}
	}
	user := User{}
	err = attributevalue.UnmarshalMap(itemMap, &user)
	if err != nil {
		return user, fmt.Errorf("Error when unpacking user: %v", err)
	}
	return user, nil
}

func (t UserTable) PutUser(u User) error {
	return putItem(t, u)
}

func (t UserTable) UpdateUser(uname string, updates map[string]interface{}) error {
	return updateItem(t, uname, updates)
}

// check that user exists and has correct password
func (t UserTable) ValidateUser(name string, password string) (bool, User, error) {
	dbUser, err := t.GetUser(name)
	if err != nil {
		if _, ok := err.(ItemNotFound); ok {
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

func (t UserTable) UserExists(name string) (bool, error) {
	_, err := t.GetUser(name)
	if err != nil {
		if _, ok := err.(ItemNotFound); ok {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
