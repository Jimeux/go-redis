package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/Jimeux/go-redis/app"
)

func main() {
	db := dynamodb.NewFromConfig(app.AWSConfig())
	createTables(db)
}

func createTables(db *dynamodb.Client) {
	if _, err := db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName:   &app.TableThreads,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: &app.PKThreads,
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: &app.PKThreads,
				KeyType:       types.KeyTypeHash,
			},
		},
	}); err != nil {
		log.Printf("%v\n", err)
	}

	if _, err := db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName:   &app.TableMessages,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: &app.PKMessages,
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: &app.SKMessages,
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: &app.PKMessages,
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: &app.SKMessages,
				KeyType:       types.KeyTypeRange,
			},
		},
	}); err != nil {
		log.Printf("%v\n", err)
	}

	if _, err := db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName:   &app.TableUsers,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: &app.PKUsers,
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: &app.PKUsers,
				KeyType:       types.KeyTypeHash,
			},
		},
	}); err != nil {
		log.Printf("%v\n", err)
	}

	if _, err := db.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName:   &app.TableReactions,
		BillingMode: types.BillingModePayPerRequest,
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: &app.PKReactions,
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: &app.SKReactions,
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: &app.PKReactions,
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: &app.SKReactions,
				KeyType:       types.KeyTypeRange,
			},
		},
	}); err != nil {
		log.Printf("%v\n", err)
	}
}
