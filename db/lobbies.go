package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Lobby struct {
	ID      string    `json:"id"`
	Host    string    `json:"host"`
	Players [4]string `json:"players"`
}

const LOBBY_TABLE_NAME = "Bird.Lobbies"

type LobbyTable struct {
	client *dynamodb.Client
}

func (t LobbyTable) Client() *dynamodb.Client {
	return t.client
}

func (t LobbyTable) Name() string {
	return LOBBY_TABLE_NAME
}

func (t LobbyTable) IndexName() string {
	return "ID"
}

func (t LobbyTable) IndexType() types.ScalarAttributeType {
	return types.ScalarAttributeTypeS
}

func MakeLobbyTable(client *dynamodb.Client) (LobbyTable, error) {
	table := LobbyTable{client}
	exists, err := tableIsInitialized(table)
	if err != nil {
		return table, fmt.Errorf("Error when checking if lobby table exists: %v", err)
	}
	if exists {
		return table, nil
	} else {
		err = initTable(table)
		return table, err
	}
}

func LobbyFromItemMap(m map[string]types.AttributeValue) (Lobby, error) {
	if m == nil {
		return Lobby{}, ItemNotFound{"Lobby"}
	}
	lobby := Lobby{}
	err := attributevalue.UnmarshalMap(m, &lobby)
	if err != nil {
		return lobby, fmt.Errorf("Error when unpacking lobby: %v", err)
	}
	return lobby, nil
}

func (t LobbyTable) GetLobby(id string) (Lobby, error) {
	itemMap, err := getItem(t, id)
	if err != nil {
		return Lobby{}, err
	}
	return LobbyFromItemMap(itemMap)
}

func (t LobbyTable) PutLobby(l Lobby) error {
	return putItem(t, l)
}

func (t LobbyTable) UpdateLobby(id string, updates map[string]interface{}) error {
	return updateItem(t, id, updates)
}

func (t LobbyTable) DeleteLobby(id string) error {
	return deleteItem(t, id)
}
