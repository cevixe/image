package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/sns"
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

	for _, record := range request.Records {

		item, err := message.FromStream(record)
		if err != nil {
			fmt.Printf("cannot read message from dynamodb stream: %v\n", err)
			continue
		}

		input, err := message.ToSNS_Input(item)
		if err != nil {
			return errors.Wrap(err, "cannot generate sns input from message")
		}

		input.TopicArn = jsii.String(h.topic)
		_, err = h.client.Publish(ctx, input)
		if err != nil {
			return errors.Wrap(err, "cannot publish message to sns topic")
		}
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
