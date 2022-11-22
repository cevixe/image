package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/cevixe/sdk/client/config"
	"github.com/cevixe/sdk/message"
	"github.com/pkg/errors"
)

type Handler struct {
	table  string
	client *dynamodb.Client
}

func (h *Handler) Handle(ctx context.Context, request events.SQSEvent) error {

	messages := make([]message.Message, 0)
	for _, record := range request.Records {

		item, err := message.FromSQS(record)
		if err != nil {
			return errors.Wrap(err, "cannot read message from sqs")
		}
		messages = append(messages, item)
	}

	err := message.Write(ctx, messages...)
	if err != nil {
		return errors.Wrap(err, "cannot read message from sqs")
	}

	return nil
}

func main() {
	table := os.Getenv("CVX_EVENT_STORE")

	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	client := dynamodb.NewFromConfig(cfg)

	handler := &Handler{
		table:  table,
		client: client,
	}

	lambda.StartWithOptions(handler.Handle, lambda.WithContext(ctx))
}
