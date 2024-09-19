package main

import (
	"context"
	"github.com/aws/aws-lambda-go/lambda"
)

func HelloWorld(ctx context.Context) (string, error) {
	return "Hello, World! - KYC processboarding - CDK", nil
}

func main() {
	lambda.Start(HelloWorld)
}
