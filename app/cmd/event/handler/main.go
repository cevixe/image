package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/jsii-runtime-go"
	"github.com/cevixe/sdk/client/config"
	"github.com/cevixe/sdk/message"
	"github.com/pkg/errors"
)

type Handler struct {
	table  string
	client *dynamodb.Client
}

func (h *Handler) Handle(ctx context.Context, request events.SQSEvent) error {

	writeItems := make([]types.TransactWriteItem, 0)
	for _, record := range request.Records {

		item, err := message.FromSQS(record)
		if err != nil {
			return errors.Wrap(err, "cannot read message from sqs")
		}

		itemMap, err := message.ToDynamodb_Map(item)
		if err != nil {
			return errors.Wrap(err, "cannot marshal message to dynamodb map")
		}

		for key, value := range itemMap {
			if value == nil {
				delete(itemMap, key)
			} else {
				switch value.(type) {
				case *types.AttributeValueMemberNULL:
					delete(itemMap, key)
				default:
					continue
				}
			}
		}

		writeItems = append(writeItems, types.TransactWriteItem{
			Put: &types.Put{
				TableName:           jsii.String(h.table),
				Item:                itemMap,
				ConditionExpression: jsii.String("attribute_not_exists(#id)"),
				ExpressionAttributeNames: map[string]string{
					"#id": "id",
				},
			},
		})
	}

	_, err := h.client.TransactWriteItems(ctx, &dynamodb.TransactWriteItemsInput{
		TransactItems: writeItems,
	})
	if err != nil {
		return errors.Wrap(err, "cannot write messages to dynamodb")
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
