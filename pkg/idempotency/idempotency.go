package idempotency

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/pkg/errors"
)

const PutCondition = "attribute_not_exists(pk)"

type IdempotencyManager interface {
	ItemAlreadyHandled(context.Context, string) (bool, error)
}

type Entry struct {
	Key         string `dynamodbav:"pk"`
	CurrentTime int64  `dynamodbav:"now"`
	ExpireAt    int64  `dynamodbav:"ttl"`
}
type Manager struct {
	table  string
	expiry time.Duration

	client *dynamodb.Client
}

func NewManager(table string, expiry time.Duration) *Manager {
	cfg, _ := config.LoadDefaultConfig(context.Background())

	return &Manager{
		table:  table,
		expiry: expiry,
		client: dynamodb.NewFromConfig(cfg),
	}
}

func (m Manager) ItemAlreadyHandled(ctx context.Context, pk string) (handled bool, err error) {
	entry, err := attributevalue.MarshalMap(Entry{
		Key:         pk,
		CurrentTime: time.Now().Unix(),
		ExpireAt:    time.Now().Add(m.expiry).Unix(),
	})
	if err != nil {
		return true, errors.Wrap(err, "unable to marshal idempotency entry")
	}

	_, err = m.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: []types.TransactWriteItem{
			{
				Put: &types.Put{
					TableName:           &m.table,
					Item:                entry,
					ConditionExpression: aws.String(PutCondition),
				},
			},
		},
	})
	if err != nil {
		handled = true
	}

	return
}
