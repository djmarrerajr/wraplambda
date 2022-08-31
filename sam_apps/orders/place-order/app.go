package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/djmarrerajr/wraplambda/pkg/idempotency"
)

const (
	oneDay = 24 * time.Hour
)

var (
	InputQueueUrl         string
	OrderTableName        string
	IdempotencyExpiryDays time.Duration

	IdempotencyManager idempotency.IdempotencyManager
)

func main() {
	err := loadEnv()
	if err != nil {
		fmt.Printf("unable to load environment: %s", err.Error())
	}

	IdempotencyManager = idempotency.NewManager(OrderTableName, IdempotencyExpiryDays)

	lambda.Start(LambdaHandler)
}

func LambdaHandler(ctx context.Context, event events.SQSEvent) (result events.SQSEventResponse, err error) {
	fmt.Printf("INPUT_QUEUE_URL: %s\n", InputQueueUrl)
	fmt.Printf("ORDER_TABLE_NAME: %s\n", OrderTableName)

	for _, record := range event.Records {
		IdempotencyManager.ItemAlreadyHandled(ctx, record.MessageId)

		time.Sleep(2 * time.Second)
		IdempotencyManager.ItemAlreadyHandled(ctx, record.MessageId)
	}

	return events.SQSEventResponse{}, nil
}

func loadEnv() (err error) {
	var exists bool
	var idempotencyExpiryDaysInt int

	InputQueueUrl, exists = os.LookupEnv("INPUT_QUEUE_URL")
	if InputQueueUrl == "" || !exists {
		return errors.New("INPUT_QUEUE_URL is missing")
	}

	OrderTableName, exists = os.LookupEnv("ORDER_TABLE_NAME")
	if OrderTableName == "" || !exists {
		return errors.New("ORDER_TABLE_NAME is missing")
	}

	idempotencyExpiryDaysStr, exists := os.LookupEnv("IDEMPOTENCY_EXPIRATION_DAYS")
	if idempotencyExpiryDaysStr == "" || !exists {
		return errors.New("IDEMPOTENCY_EXPIRATION_DAYS is missing")
	}

	idempotencyExpiryDaysInt, err = strconv.Atoi(idempotencyExpiryDaysStr)
	if err != nil {
		return errors.New("IDEMPOTENCY_EXPIRATION_DAYS is invalid")
	}

	IdempotencyExpiryDays = time.Duration(oneDay * time.Duration(idempotencyExpiryDaysInt))

	return nil
}
