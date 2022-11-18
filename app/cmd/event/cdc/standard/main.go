package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/sdk/event"
)

type Handler struct {
	topic  string
	client *sns.Client
}

func (h *Handler) Handle(ctx context.Context, request events.DynamoDBEvent) error {

	entries := make([]types.PublishBatchRequestEntry, 0)
	for _, record := range request.Records {
		item, err := event.From_DynamoDBEventRecord(record)
		if err != nil {
			return err
		}

		entries = append(entries, event.To_SNSPublishBatchRequestEntry(item))
	}

	if len(entries) == 0 {
		return nil
	}

	_, err := h.client.PublishBatch(ctx, &sns.PublishBatchInput{
		TopicArn:                   jsii.String(h.topic),
		PublishBatchRequestEntries: entries,
	})
	return err
}

func main() {
	region := os.Getenv("AWS_REGION")
	topic := os.Getenv("CVX_EVENT_BUS")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	client := sns.NewFromConfig(cfg)

	handler := &Handler{
		topic:  topic,
		client: client,
	}

	lambda.Start(handler.Handle)
}
