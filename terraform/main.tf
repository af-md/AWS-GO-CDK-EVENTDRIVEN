// output and variables could be added


terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
  }
}
# Configure the AWS Provider
provider "aws" {
  region = "eu-west-1"
  shared_credentials_files = ["/Users/afzal.akmal.muhammad/.aws/credentials"]
}