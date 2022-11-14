package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/sdk/event"
)

type Handler struct {
	topic  string
	client *sns.Client
}

func (h *Handler) Handle(ctx context.Context, request events.DynamoDBEvent) error {

	for _, record := range request.Records {
		item := event.From_DynamoDBEventRecord(record)
		input := event.To_SNSPublishInput(item)
		input.TopicArn = jsii.String(h.topic)
		_, err := h.client.Publish(ctx, &input)
		return err
	}

	return nil
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
