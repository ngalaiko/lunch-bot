package main

import (
	"context"

	"lunch/pkg/handler"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.StartWithContext(context.Background(), handler.Handle)
}
