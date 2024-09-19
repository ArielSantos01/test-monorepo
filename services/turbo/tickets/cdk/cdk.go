package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"os"
	"stackcdk"
)

func main() {

	defer jsii.Close()

	app := awscdk.NewApp(nil)

	_, err := stackcdk.NewCdkStack(app, "turbo-tickets-test-full", &stackcdk.CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	if err != nil {
		panic(err)
	}

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT")),
		Region:  jsii.String("us-east-1"),
	}
}
