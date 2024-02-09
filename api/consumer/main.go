package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func init() {
}

func lambdaHandler(context context.Context, event events.EventBridgeEvent) {
	fmt.Println("Hello this is the lambda")
	fmt.Println(event)
}

func main() {
	lambda.Start(lambdaHandler)
}
