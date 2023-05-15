package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handle(ctx context.Context, req json.RawMessage) ([]byte, error) {
	log.Printf("%s", req)
	return nil, nil
}

func main() {
	lambda.Start(handle)
}
