package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/cevixe/sdk/event"
	"github.com/pkg/errors"
)

type Handler struct {
	table  string
	client *dynamodb.Client
}

func (h *Handler) Handle(ctx context.Context, request events.SQSEvent) error {

	writes := make([]types.WriteRequest, 0)
	for _, record := range request.Records {
		item := event.From_SQSMessage(record)
		statement, err := event.To_DynamodbWriteRequest(item)
		if err != nil {
			return errors.Wrap(err, "cannot generate write request")
		}
		writes = append(writes, *statement)
	}

	_, err := h.client.BatchWriteItem(ctx, &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			h.table: writes,
		},
	})
	return err
}

func main() {
	region := os.Getenv("AWS_REGION")
	table := os.Getenv("CVX_EVENT_STORE")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	handler := &Handler{
		table:  table,
		client: client,
	}

	lambda.Start(handler.Handle)
}
