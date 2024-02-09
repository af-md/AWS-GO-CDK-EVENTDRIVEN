package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/jsii-runtime-go"

	integ "github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2integrations"
	"github.com/aws/constructs-go/constructs/v10"
)

type HttpApiEventbridgeStackProps struct {
	awscdk.StackProps
}

const region = "eu-west-2"

func NewHttpApiEventbridgeStack(scope constructs.Construct, id string, props *HttpApiEventbridgeStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// create EventBridge event bus
	eventBus := awsevents.NewEventBus(stack, jsii.String("myEventBus"), &awsevents.EventBusProps{
		EventBusName: jsii.String("MyEventBus"),
	})

	// create EventBridge rule
	eventLoggerRule := awsevents.NewRule(stack, jsii.String("myEventLoggerRule"), &awsevents.RuleProps{
		Description: jsii.String("Log all events"),
		EventBus:    eventBus,
		EventPattern: &awsevents.EventPattern{
			Region: jsii.Strings(*props.Env.Region),
		},
	})

	// add CloudWatch log group as target
	logGroup := awslogs.NewLogGroup(stack, jsii.String("MyEventLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/aws/events/MyEventBus"),
	})

	// add target to rule
	eventLoggerRule.AddTarget(awseventstargets.NewCloudWatchLogGroup(logGroup, &awseventstargets.LogGroupProps{}))

	// create HTTP API
	httpApi := awsapigatewayv2.NewHttpApi(stack, jsii.String("myHttpApi"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("myHttpApi"),
	})

	// create IAM role for API Gateway to put events to EventBridge
	apiRole := awsiam.NewRole(stack, jsii.String("myEventBridgeIntegrationRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("apigateway.amazonaws.com"), &awsiam.ServicePrincipalOpts{}),
	})

	// add policy to role
	apiRole.AddToPolicy(
		awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
			Effect:    awsiam.Effect_ALLOW,
			Resources: jsii.Strings(*eventBus.EventBusArn()),
			Actions:   jsii.Strings("events:putEvents"),
		}),
	)

	// create HTTP API integration
	eventbridgeIntegration := awsapigatewayv2.NewCfnIntegration(stack, jsii.String("myEventbridgeIntegration"), &awsapigatewayv2.CfnIntegrationProps{
		ApiId:              httpApi.HttpApiId(),
		IntegrationType:    jsii.String("AWS_PROXY"),
		IntegrationSubtype: jsii.String("EventBridge-PutEvents"),
		CredentialsArn:     apiRole.RoleArn(),
		RequestParameters: &map[string]string{
			"Source":       "WebApp",
			"DetailType":   "MyDetailType",
			"Detail":       "$request.body",
			"EventBusName": *eventBus.EventBusArn(),
		},
		PayloadFormatVersion: jsii.String("1.0"),
		TimeoutInMillis:      jsii.Number(10000.0),
	})

	// create HTTP API route
	awsapigatewayv2.NewCfnRoute(stack, jsii.String("myEventRoute"), &awsapigatewayv2.CfnRouteProps{
		ApiId:    httpApi.HttpApiId(),
		RouteKey: jsii.String("POST /"),
		Target:   jsii.String("integrations/" + *eventbridgeIntegration.Ref()),
	})

	// log HTTP API endpoint URL
	awscdk.NewCfnOutput(stack, jsii.String("apiUrl"), &awscdk.CfnOutputProps{
		Value:       httpApi.Url(),
		Description: jsii.String("HTTP API endpoint URL"),
	})

	// create a new lambda, add it to the api gateway endpoint, deploy it, make the lambda output something. ultimately we would like this to store the events inside dynamo
	// create Lambda function

	getHandler := awslambda.NewFunction(stack, jsii.String("myGoHandler"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_GO_1_X(),
		Code:awslambda.Code_FromAsset(jsii.String("../lambda-handler"), &awss3assets.AssetOptions{}),
		Handler: jsii.String("main"),
	})

	// create HTTP API
	httpApi2 := awsapigatewayv2.NewHttpApi(stack, jsii.String("myHttpApi2"), &awsapigatewayv2.HttpApiProps{
		ApiName: jsii.String("myHttpApi2"),
	})
	// add route to HTTP API
	httpApi2.AddRoutes(&awsapigatewayv2.AddRoutesOptions{
		Path:        jsii.String("/"),
		Methods:     &[]awsapigatewayv2.HttpMethod{awsapigatewayv2.HttpMethod_GET},
		Integration: integ.NewHttpLambdaIntegration(jsii.String("MyHttpLambdaIntegration"), getHandler, &integ.HttpLambdaIntegrationProps{}),
	})

	// log HTTP API endpoint
	awscdk.NewCfnOutput(stack, jsii.String("myHttpApiEndpoint"), &awscdk.CfnOutputProps{
		Value:       httpApi2.ApiEndpoint(),
		Description: jsii.String("HTTP API Endpoint"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewHttpApiEventbridgeStack(app, "HttpApiEventbridgeStack", &HttpApiEventbridgeStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region: jsii.String(getRegion()),
	}
}

func getRegion() string {
	envReg := os.Getenv("AWS_DEFAULT_REGION")
	if envReg != "" {
		return envReg
	}
	return region
}