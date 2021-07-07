package app

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	exp "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/oklog/ulid/v2"
)

type Repository struct {
	db *dynamodb.Client
}

func NewRepository(db *dynamodb.Client) *Repository {
	return &Repository{db: db}
}

func (r *Repository) SaveThread(ctx context.Context, title string) (*Thread, error) {
	t := NewThread(generateID(), title)
	av, err := attributevalue.MarshalMap(t)
	if err != nil {
		return nil, err
	}

	if _, err := r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: &TableThreads,
	}); err != nil {
		return nil, err
	}
	return t, nil
}

func (r *Repository) SaveMessage(ctx context.Context, threadID, content string, uid int64) (*Message, error) {
	m := NewMessage(threadID, content, uid)
	av, err := attributevalue.MarshalMap(m)
	if err != nil {
		return nil, err
	}

	if _, err := r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: &TableMessages,
	}); err != nil {
		return nil, err
	}
	return m, nil
}

func (r *Repository) FindMessages(ctx context.Context, threadID string, min, max, limit int64) ([]*Message, error) {
	expr, _ := exp.NewBuilder().
		WithKeyCondition(exp.Key(PKMessages).Equal(exp.Value(threadID)).
			And(exp.Key(SKMessages).Between(exp.Value(min), exp.Value(max)))).
		Build()

	res, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &TableMessages,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		Limit:                     aws.Int32(int32(limit)),
	})
	if err != nil {
		return nil, err
	}

	var msgs []*Message
	if err := attributevalue.UnmarshalListOfMaps(res.Items, &msgs); err != nil {
		return nil, err
	}
	return msgs, nil
}

func (r *Repository) SaveReaction(ctx context.Context, threadID string, sentAt int64, userID int64, kind string) (*Reaction, error) {
	rct := NewReaction(threadID, sentAt, userID, kind)
	av, err := attributevalue.MarshalMap(rct)
	if err != nil {
		return nil, err
	}

	if _, err := r.db.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      av,
		TableName: &TableReactions,
	}); err != nil {
		return nil, err
	}
	return rct, nil
}

func (r *Repository) FindReactions(ctx context.Context, threadID string, min, max int64) ([]*Reaction, error) {
	minVal := strconv.FormatInt(min, 10)
	maxVal := strconv.FormatInt(max+1, 10) // Add 1 for LTE
	expr, _ := exp.NewBuilder().
		WithKeyCondition(exp.Key(PKReactions).Equal(exp.Value(threadID)).
			And(exp.Key(SKReactions).Between(exp.Value(minVal), exp.Value(maxVal)))).
		Build()

	res, err := r.db.Query(ctx, &dynamodb.QueryInput{
		TableName:                 &TableReactions,
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	var rr []*Reaction
	if err := attributevalue.UnmarshalListOfMaps(res.Items, &rr); err != nil {
		return nil, err
	}
	return rr, nil
}

func generateID() string {
	t := time.Now().UTC()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
