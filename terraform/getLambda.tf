
//terraform plan
// terraform apply


// ALL INFRA IS DEPLOYED INTO EU-WEST-1 IRELAND


# See also the following AWS managed policy: AWSLambdaBasicExecutionRole

// permissions: create a policy document, attach to the role (lambda assume role), attach the role to the lambda
data "aws_iam_policy_document" "getLambda_policy_document" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents",

      "dynamodb:BatchGetItem",
      "dynamodb:GetItem",
      "dynamodb:Query",
      "dynamodb:Scan",
      "dynamodb:BatchWriteItem",
      "dynamodb:PutItem",
      "dynamodb:UpdateItem"
    ]
  // wildcards to access any logs or dynamodb
    resources = [
       "arn:aws:logs:*:*:*",
       "arn:aws:dynamodb:*:*:*",
      ]
  }
}

resource "aws_iam_policy" "getLambda_policy" {
  name        = "getLambda"
  path        = "/"
  description = "IAM policy for logging from a lambda"
  policy      = data.aws_iam_policy_document.getLambda_policy_document.json
}

resource "aws_iam_role_policy_attachment" "getLambda" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.getLambda_policy.arn
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}

resource "aws_iam_role" "iam_for_getLambda" {
  name               = "iam_for_getLambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "archive_file" "getlambda" {
  type        = "zip"
  source_file = "../api/lambda-handler/lambdaHandler"
  output_path = "../api/lambda-handler/lambdaHandler.zip"
}

resource "aws_lambda_function" "getHandler" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = data.archive_file.getlambda.output_path
  function_name = "getLambda"
  role          = aws_iam_role.iam_for_getLambda.arn
  handler       = "lambdaHandler"
  source_code_hash = data.archive_file.lambda.output_base64sha256
  runtime = "go1.x"
}

resource "aws_lambda_function_url" "addLambdaUrl" {
  function_name      = aws_lambda_function.addHandler.function_name
  authorization_type = "NONE"
}

## output of the lambda. what do i want ? for add i would like for it to output the apigateway or a function url 
output "lambdaOutput" {
    value = aws_lambda_function_url.addLambdaUrl.function_url
}