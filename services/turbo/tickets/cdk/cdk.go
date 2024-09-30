package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/jsii-runtime-go"
	"os"
	"path/filepath"
	"pkl"
	"stackcdk"
)

var pklEnvs map[string]any

func loadPklConfig() (map[string]any, error) {
	workdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	servicePath := filepath.Join(workdir, "..")

	isLocal := false
	if os.Getenv("LOCAL") == "true" {
		isLocal = true
	}
	return pkl.Pkl(servicePath, isLocal)

}

func init() {
	vals, err := loadPklConfig()
	if err != nil {
		panic(err)
	}
	pklEnvs = vals // pklEnvs es un map con las variables de pkl
}

func main() {

	defer jsii.Close()

	app := awscdk.NewApp(nil)

	_, err := stackcdk.NewCdkStack(app, "turbo-tickets-test-full", &stackcdk.CdkStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	},
		stackcdk.WithPklConfig(pklEnvs),
	)

	if err != nil {
		panic(err)
	}

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT")),
		Region:  jsii.String(os.Getenv("AWS_DEFAULT_REGION")),
	}
}
