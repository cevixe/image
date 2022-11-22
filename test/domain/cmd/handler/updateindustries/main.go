package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cevixe/sdk/message"
	"github.com/cevixe/sdk/result"
	"github.com/cevixe/sdk/runtime"
)

func handle(ctx context.Context, msg message.Message) (result.Result, error) {

	jsonString, _ := json.Marshal(msg)
	fmt.Println(jsonString)
	return nil, nil
}

func main() {
	runtime.Start(handle)
}
