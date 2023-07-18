package db

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/quevivasbien/bird-backend/game"
)

const BID_TABLE_NAME = "Bird.Bids"

type BidTable struct {
	client *dynamodb.Client
}

func (t BidTable) Client() *dynamodb.Client {
	return t.client
}

func (t BidTable) Name() string {
	return BID_TABLE_NAME
}

func (t BidTable) IndexName() string {
	return "ID"
}

func (t BidTable) IndexType() types.ScalarAttributeType {
	return types.ScalarAttributeTypeS
}

func MakeBidTable(client *dynamodb.Client) (BidTable, error) {
	table := BidTable{client}
	exists, err := tableIsInitialized(table)
	if err != nil {
		return table, fmt.Errorf("Error when checking if bid table exists: %v", err)
	}
	if exists {
		return table, nil
	} else {
		err = initTable(table)
		return table, err
	}
}

func (t BidTable) GetBidState(id string) (game.BidState, error) {
	itemMap, err := getItem(t, id)
	if err != nil {
		return game.BidState{}, err
	}
	if itemMap == nil {
		return game.BidState{}, ItemNotFound{"BidState"}
	}
	lobby := game.BidState{}
	err = attributevalue.UnmarshalMap(itemMap, &lobby)
	if err != nil {
		return lobby, fmt.Errorf("Error when unpacking lobby: %v", err)
	}
	return lobby, nil
}

func (t BidTable) PutBidState(b game.BidState) error {
	return putItem(t, b)
}

func (t BidTable) UpdateLobby(id string, updates map[string]interface{}) error {
	return updateItem(t, id, updates)
}
