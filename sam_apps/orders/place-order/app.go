package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/djmarrerajr/wraplambda/pkg/idempotency"
	"github.com/djmarrerajr/wraplambda/pkg/utils"
)

const (
	pkValue = "test-key"
)

var (
	InputQueueUrl  string
	OrderTableName string
)

func main() {
	err := loadEnv()
	if err != nil {
		fmt.Printf("unable to load environment: %s", err.Error())
	}

	lambda.Start(idempotency.AtMostOnce(LambdaHandler, OrderTableName))
}

func LambdaHandler(ctx context.Context, event events.SQSEvent) (result events.SQSEventResponse, err error) {
	defer utils.HandlePanic(func(panicMsg interface{}) {
		err = fmt.Errorf("panic encountered: %s", panicMsg)
	})

	fmt.Printf("INPUT_QUEUE_URL: %s\n", InputQueueUrl)
	fmt.Printf("ORDER_TABLE_NAME: %s\n", OrderTableName)

	return events.SQSEventResponse{}, nil
}

func loadEnv() error {
	var exists bool

	InputQueueUrl, exists = os.LookupEnv("INPUT_QUEUE_URL")
	if InputQueueUrl == "" || !exists {
		return errors.New("INPUT_QUEUE_URL is missing")
	}

	OrderTableName, exists = os.LookupEnv("ORDER_TABLE_NAME")
	if OrderTableName == "" || !exists {
		return errors.New("ORDER_TABLE_NAME is missing")
	}

	return nil
}
