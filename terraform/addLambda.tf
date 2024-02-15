//terraform plan
// terraform apply

# See also the following AWS managed policy: AWSLambdaBasicExecutionRole

// permissions: create a policy document, attach to the role (lambda assume role), attach the role to the lambda
data "aws_iam_policy_document" "addLambda_policy_document" {
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

resource "aws_iam_policy" "addLambda_policy" {
  name        = "addLambda"
  path        = "/"
  description = "IAM policy for logging from a lambda"
  policy      = data.aws_iam_policy_document.addLambda_policy_document.json
}

resource "aws_iam_role_policy_attachment" "addLambda" {
  role       = aws_iam_role.iam_for_lambda.name
  policy_arn = aws_iam_policy.addLambda_policy.arn
}

resource "aws_iam_role" "iam_for_lambda" {
  name               = "iam_for_lambda"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

data "archive_file" "lambda" {
  type        = "zip"
  source_file = "../api/consumer/lambdaHandler"
  output_path = "../api/consumer/lambdaHandler.zip"
}

resource "aws_lambda_function" "addHandler" {
  # If the file is not in the current working directory you will need to include a
  # path.module in the filename.
  filename      = data.archive_file.lambda.output_path
  function_name = "addLambda"
  role          = aws_iam_role.iam_for_lambda.arn
  handler       = "lambdaHandler"

  source_code_hash = data.archive_file.lambda.output_base64sha256

  runtime = "go1.x"
}

