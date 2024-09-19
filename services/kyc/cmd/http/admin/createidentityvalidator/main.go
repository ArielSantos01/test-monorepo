package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

func HelloWorld(ctx context.Context) (string, error) {
	return "Hello, World! - KYC createidentityvalidator - CDK", nil
}

func main() {
	lambda.Start(HelloWorld)
}
