package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/cevixe/sdk/common/json"
	"github.com/cevixe/sdk/event"
)

type Handler struct {
}

func (h *Handler) Handle(ctx context.Context, request events.SQSMessage) error {

	event := event.From_SQSMessage(request)
	jsonString := json.Marshal(event)
	fmt.Println(jsonString)
	return nil
}

func main() {
	region := os.Getenv("AWS_REGION")

	_, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
	)

	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	handler := &Handler{}

	lambda.Start(handler.Handle)
}
