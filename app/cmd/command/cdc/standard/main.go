package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/sdk/client/config"
	"github.com/cevixe/sdk/message"
	"github.com/pkg/errors"
)

type Handler struct {
	topic  string
	client *sns.Client
}

func (h *Handler) Handle(ctx context.Context, request events.DynamoDBEvent) error {

	entries := make([]types.PublishBatchRequestEntry, 0)
	for _, record := range request.Records {

		item, err := message.FromStream(record)
		if err != nil {
			return errors.Wrap(err, "cannot read message from dynamodb stream")
		}

		entry, err := message.ToSNS_Entry(item)
		if err != nil {
			return errors.Wrap(err, "cannot generate sns input from message")
		}

		entries = append(entries, *entry)
	}

	if len(entries) == 0 {
		return nil
	}

	_, err := h.client.PublishBatch(ctx, &sns.PublishBatchInput{
		TopicArn:                   jsii.String(h.topic),
		PublishBatchRequestEntries: entries,
	})
	if err != nil {
		return errors.Wrap(err, "cannot publish messages to sns topic")
	}

	return nil
}

func main() {
	topic := os.Getenv("CVX_EVENT_BUS")

	ctx := context.Background()
	cfg := config.NewConfig(ctx)
	client := sns.NewFromConfig(cfg)

	handler := &Handler{
		topic:  topic,
		client: client,
	}

	lambda.StartWithOptions(handler.Handle, lambda.WithContext(ctx))
}
