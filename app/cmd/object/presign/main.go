package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/cevixe/sdk/client/config"
	"github.com/cevixe/sdk/object"
)

type Operation string

const (
	Operation_Upload   Operation = "upload"
	Operation_Download Operation = "download"
)

type Input struct {
	Operation Operation `field:"required" json:"operation"`
	Name      string    `field:"required" json:"name"`
}

type Output struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Handler struct {
	client object.Client
}

func (h *Handler) Handle(ctx context.Context, input *Input) (*Output, error) {
	switch input.Operation {
	case Operation_Upload:
		url, err := h.client.UploadURL(ctx, input.Name, 5*time.Minute)
		if err != nil {
			return nil, err
		}
		return &Output{Name: input.Name, URL: *url}, nil
	case Operation_Download:
		exists, err := h.client.Exists(ctx, input.Name)
		if err != nil {
			return nil, err
		}
		if *exists {
			return nil, nil
		}
		url, err := h.client.DownloadURL(ctx, input.Name, 5*time.Minute)
		if err != nil {
			return nil, err
		}
		return &Output{Name: input.Name, URL: *url}, nil
	default:
		return nil, fmt.Errorf("unknown link generator operation")
	}
}

func main() {
	bucket := os.Getenv("CVX_OBJECT_STORE")

	ctx := context.Background()
	cfg := config.NewConfig(ctx)

	standardClient := s3.NewFromConfig(cfg)
	presignClient := s3.NewPresignClient(standardClient)

	objectClient := object.NewClient(bucket, presignClient, standardClient)
	handler := &Handler{client: objectClient}

	lambda.StartWithOptions(handler.Handle, lambda.WithContext(ctx))
}
