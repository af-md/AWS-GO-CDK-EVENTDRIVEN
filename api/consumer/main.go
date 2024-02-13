package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/jsii-runtime-go"
)

func init() {
}

func lambdaHandler(context context.Context, event events.EventBridgeEvent) {

	log, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("could not create log: %v", err))
	}

	eventjs, err := json.Marshal(event)

	if err != nil {
		log.Fatal("could not load AWS config", zap.Error(err))
	}

	log.Info(string(eventjs))

	cfg, err := config.LoadDefaultConfig(context, config.WithRegion("eu-west-2"))
	if err != nil {
		log.Fatal("unable to load SDK config, %v", zap.Error(err))
	}

	if err != nil {
		log.Fatal("unable to marshal event into dynamodb attribute %v", zap.Error(err))
	}

	itemV := makeItem(event)

	table := "EventTable"
	input := &dynamodb.PutItemInput{
		TableName: jsii.String(table),
		Item:      itemV,
	}

	// Using the Config value,sl create the DynamoDB client
	dc := dynamodb.NewFromConfig(cfg)
	_, err = dc.PutItem(context, input)
	if err != nil {
		log.Fatal("unable to put item into dynamodb %v", zap.Error(err))
	} else {
		log.Info("successfully added the item inside the event")
	}

}

func main() {
	lambda.Start(lambdaHandler)
}

func makeItem(event events.EventBridgeEvent) map[string]types.AttributeValue {
	item := map[string]types.AttributeValue{
		"id": &types.AttributeValueMemberS{
			Value: event.ID,
		},
		"source": &types.AttributeValueMemberS{
			Value: event.Source,
		},
	}
	return item
}
