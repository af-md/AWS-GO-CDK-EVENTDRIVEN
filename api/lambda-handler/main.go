package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/jsii-runtime-go"
	"go.uber.org/zap"
)

type response struct {
	Source string `json:"source"`
	Id     string `json:"id"`
}

func init() {
}

func handleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	log, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("could not create log: %v", err))
	}

	if err != nil {
		log.Fatal("could not load AWS config", zap.Error(err))
	}

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("eu-west-2"))
	if err != nil {
		log.Fatal("unable to load SDK config, %v", zap.Error(err))
	}

	dc := dynamodb.NewFromConfig(cfg)

	items, err := dc.Scan(ctx, &dynamodb.ScanInput{
		TableName: jsii.String("EventTable"),
	})

	if err != nil {
		log.Fatal("unable to scan from dynamo, %v", zap.Error(err))
	}

	var evs []response

	err = attributevalue.UnmarshalListOfMaps(items.Items, &evs)
	if err != nil {
		log.Fatal("unable to unmarshal list of maps items, %v", zap.Error(err))
	}

	body, err := json.Marshal(evs)
	if err != nil {
		return events.APIGatewayProxyResponse{Body: string("Error parsing payload"), StatusCode: 400}, err
	}
	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 200}, nil
}

func main() {
	lambda.Start(handleRequest)
}
