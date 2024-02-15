// create an eventbridge bus, rule, target

resource "aws_cloudwatch_event_bus" "messenger" {
  name = "events-finance"
}

resource "aws_cloudwatch_event_rule" "addEvent" {
  name        = "addEvent"
  event_pattern = jsonencode({
    source = [
      "kafka"
    ]
  })
}

resource "aws_cloudwatch_event_target" "lambda" {
  rule      = aws_cloudwatch_event_rule.addEvent.name
  target_id = "sendtolambda"
  arn       = aws_lambda_function.addHandler.arn
}


## integration with api gateway


resource "aws_iam_role" "ApiGatewayEventBridgeRole" {
  name = "ApiGatewayEventBridgeRole"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Action = "sts:AssumeRole"
      Effect = "Allow"
      Sid    = ""
      Principal = {
        Service = "apigateway.amazonaws.com"
      }
      }
    ]
  })
}

resource "aws_iam_policy" "EBPutEvents" {
  name = "EBPutEvents"

  policy = <<POLICY
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Action": "events:PutEvents",
            "Resource": "arn:aws:events:eu-west-1:109678184122:event-bus/events-finance"
        }
    ]
}
POLICY
}


resource "aws_iam_role_policy_attachment" "apigwy_policy" {
  role       = aws_iam_role.ApiGatewayEventBridgeRole.name
  policy_arn = aws_iam_policy.EBPutEvents.arn
}



resource "aws_api_gateway_rest_api" "my_api" {
  name        = "MyAPI"
  description = "My API Gateway"
 
}

resource "aws_api_gateway_resource" "root_resource" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  parent_id   = aws_api_gateway_rest_api.my_api.root_resource_id
  path_part   = "resource"
}

resource "aws_api_gateway_method" "post_method" {
  rest_api_id   = aws_api_gateway_rest_api.my_api.id
  resource_id   = aws_api_gateway_resource.root_resource.id
  http_method   = "POST"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "eventbridge_integration" {
  rest_api_id             = aws_api_gateway_rest_api.my_api.id
  resource_id             = aws_api_gateway_resource.root_resource.id
  http_method             = aws_api_gateway_method.post_method.http_method
  integration_http_method = "POST"
  type                    = "AWS"
  uri                     = "arn:aws:apigateway:eu-west-1:events:action/PutEvents"
  credentials          = aws_iam_role.ApiGatewayEventBridgeRole.arn
}

resource "aws_api_gateway_method_response" "post_method_response" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.root_resource.id
  http_method = aws_api_gateway_method.post_method.http_method
  status_code = 200
}

resource "aws_api_gateway_integration_response" "eventbridge_integration_response" {
  rest_api_id = aws_api_gateway_rest_api.my_api.id
  resource_id = aws_api_gateway_resource.root_resource.id
  http_method = aws_api_gateway_method.post_method.http_method
  status_code = aws_api_gateway_method_response.post_method_response.status_code
  response_templates = {
    "application/json" = ""
  }
}

output "eventbus_url" {
  value = aws_api_gateway_integration.eventbridge_integration.
}